package repositories

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/logger"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	uuid "github.com/satori/go.uuid"
)

type ResponseDTO struct {
	Code        int
	Message     string
	MessageCode string
	Status      bool
	Error       error
}

type QrPrefix struct {
	QrPrefix string
	Qr       []string
	Count    int
}

// constructQrPrefix : construct qr prefix
func constructQrPrefix(qr []string) []QrPrefix {
	mapQrPrefix := make(map[string]int, 0)
	mapQr := make(map[string][]string, 0)
	for _, item := range qr {
		qrPrefix := strings.Split(item, "-")
		mapQrPrefix[fmt.Sprintf("%s-%s-", qrPrefix[0], qrPrefix[1])]++
		mapQr[fmt.Sprintf("%s-%s-", qrPrefix[0], qrPrefix[1])] = append(mapQr[fmt.Sprintf("%s-%s-", qrPrefix[0], qrPrefix[1])], item)
	}

	qrPrefixes := make([]QrPrefix, 0)
	for k, v := range mapQrPrefix {
		qrPrefixes = append(qrPrefixes, QrPrefix{
			QrPrefix: k,
			Qr:       mapQr[k],
			Count:    v,
		})
	}

	return qrPrefixes
}

// RedeemTicket : redeem ticket
func RedeemTicketV2(userData *entities.Users, req *requests.RedeemReqV2) (map[string]interface{}, int, string, string, bool) {
	resp := make(map[string]interface{}, 0)

	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return resp, http.StatusBadRequest, err.Error(), "TRANSACTION_OTA_DATE_INVALID", false
	}

	today := time.Now().Truncate(24 * time.Hour)
	if visitDate.Before(today) {
		return resp, http.StatusBadRequest, "Visit date must be today or later", "TRANSACTION_OTA_DATE_INVALID", false
	}

	sort.Strings(req.QR)
	batchSize := 10
	expectedResponses := len(req.QR) / batchSize
	if len(req.QR)%batchSize != 0 {
		expectedResponses++
	}

	wg := new(sync.WaitGroup)
	chResp := make(chan ResponseDTO, expectedResponses)
	for start, end := 0, 0; start < len(req.QR); start = end {
		end = start + batchSize
		if end > len(req.QR) {
			end = len(req.QR)
		}

		batch := req.QR[start:end]

		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()

			// checking if qr is not redeemed
			var checkOtaInventoryDetails []entities.OtaInventoryDetail

			if err := db.DB[0].Raw(`
			SELECT oid2.*
			FROM ota_inventory_detail oid2
			JOIN ota_inventory oi ON oi.id = oid2.ota_inventory_id
			WHERE oi.agent_id = ?
			AND oid2.qr IN (?)
			AND oid2.redeem_date IS NOT NULL
			LIMIT ?`, userData.Typeid, batch, len(batch)).
				Scan(&checkOtaInventoryDetails).Error; err != nil {
				chResp <- ResponseDTO{
					Code:        http.StatusInternalServerError,
					Message:     err.Error(),
					MessageCode: "TRANSACTION_OTA_FAILED",
					Status:      false,
					Error:       err,
				}
				return
			}

			if len(checkOtaInventoryDetails) > 0 {
				chResp <- ResponseDTO{
					Code:        http.StatusBadRequest,
					Message:     fmt.Sprintf("QR %s has been redeemed", checkOtaInventoryDetails[0].QR),
					MessageCode: "TRANSACTION_OTA_REDEEMED",
					Status:      false,
					Error:       errors.New("qr has been redeemed"),
				}
				return
			}
			// end checking if qr is not redeemed

			// start database transaction
			tx := db.DB[0].Begin()
			defer func() {
				if r := recover(); r != nil {
					logger.Warning("Database Rollback", "400", false, fmt.Sprintf("%+v", r))
					tx.Rollback()
				}
			}()

			if err := tx.Error; err != nil {
				chResp <- ResponseDTO{
					Code:        http.StatusInternalServerError,
					Message:     err.Error(),
					MessageCode: "TRANSACTION_OTA_FAILED",
					Status:      false,
					Error:       err,
				}
				return
			}

			qrPrefixes := constructQrPrefix(batch)
			for _, qrPrefix := range qrPrefixes {
				// start fetching ota inventory details by given qr
				var otaInventoryDetails []entities.OtaInventoryDetail

				if err := db.DB[0].Raw(`
				SELECT oid2.*
				FROM ota_inventory_detail oid2
				JOIN ota_inventory oi ON oi.id = oid2.ota_inventory_id
				WHERE oi.agent_id = ?
				AND oid2.redeem_date IS NULL
				AND oid2.void_date IS NULL
				AND (
					(oid2.qr IN (?))
					OR
					((oid2.qr IS NULL OR oid2.qr = '') AND oid2.qr_prefix = ?)
				)
				LIMIT ?`, userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count).
					Scan(&otaInventoryDetails).Error; err != nil {
					chResp <- ResponseDTO{
						Code:        http.StatusInternalServerError,
						Message:     err.Error(),
						MessageCode: "TRANSACTION_OTA_FAILED",
						Status:      false,
						Error:       err,
					}
					return
				}

				if len(otaInventoryDetails) == 0 {
					chResp <- ResponseDTO{
						Code:        http.StatusNotFound,
						Message:     "QR not found",
						MessageCode: "TRANSACTION_OTA_NOT_FOUND",
						Status:      false,
						Error:       errors.New("qr not found"),
					}
					return
				}

				if len(otaInventoryDetails) < qrPrefix.Count {
					chResp <- ResponseDTO{
						Code:        http.StatusBadRequest,
						Message:     "QR maximum exceeded",
						MessageCode: "TRANSACTION_OTA_MAX_EXCEEDED",
						Status:      false,
						Error:       errors.New("qr maximum exceeded"),
					}
					return
				}
				// end fetching ota inventory details by given qr

				// start filtering mid
				mids := make([]string, 0)
				for _, item := range otaInventoryDetails {
					if item.ExpiryDate.Before(time.Now()) {
						qr := item.QR
						if qr == "" {
							qr = item.QrPrefix
						}

						chResp <- ResponseDTO{
							Code:        http.StatusBadRequest,
							Message:     fmt.Sprintf("Ticket %s has expired", qr),
							MessageCode: "TRANSACTION_OTA_EXPIRED",
							Status:      false,
							Error:       errors.New("ticket has expired"),
						}
						return
					}

					mids = append(mids, item.GroupMid)
				}

				filteredMids := helper.RemoveDuplicateStr(mids)
				// end filtering mid

				// start main process
				for _, mid := range filteredMids {
					// start initializing variables
					now := time.Now()
					newTickID := uuid.NewV4()
					tickAmount := new(float32)
					qrData := make([]entities.OtaInventoryDetail, 0)
					// end initializing variables

					// start updating and appending qr data by given mid
					for i, otaInventoryDetail := range otaInventoryDetails {
						if otaInventoryDetail.GroupMid == mid {
							// start checking & updating ota inventory detail
							if otaInventoryDetail.QrPrefix != "" {
								otaInventoryDetail.QR = qrPrefix.Qr[i]
							}

							otaInventoryDetail.RedeemDate = &now

							if err := tx.Save(&otaInventoryDetail).Error; err != nil {
								chResp <- ResponseDTO{
									Code:        http.StatusInternalServerError,
									Message:     err.Error(),
									MessageCode: "TRANSACTION_OTA_FAILED",
									Status:      false,
									Error:       err,
								}
								return
							}

							qrData = append(qrData, otaInventoryDetail)
							*tickAmount += otaInventoryDetail.TrfAmount
							// end checking & updating ota inventory detail
						}
					}
					// end updating and appending qr data by given mid

					// start initialize new ticket
					tickStan := now.UnixNano()
					microStan := tickStan / (int64(time.Millisecond) / int64(time.Nanosecond))
					tickNumber := fmt.Sprintf("TWC.5.%d.%d", userData.Typeid, microStan)
					newTicket := entities.TickModel{
						Tick_id:             newTickID,
						Tick_stan:           int(tickStan),
						Tick_number:         tickNumber,
						Tick_mid:            mid,
						Tick_src_type:       5,
						Tick_src_id:         fmt.Sprintf("%d", userData.Typeid),
						Tick_src_inv_num:    req.OtaOrderID,
						Tick_amount:         *tickAmount,
						Tick_emoney:         0,
						Tick_purc:           now.Format("2006-01-02 15:04:05"),
						Tick_issuing:        now.Format("2006-01-02 15:04:05"),
						Tick_date:           req.VisitDate,
						Tick_total_payment:  *tickAmount,
						Tick_payment_method: "OTA",
					}

					if err := tx.Create(&newTicket).Error; err != nil {
						chResp <- ResponseDTO{
							Code:        http.StatusInternalServerError,
							Message:     err.Error(),
							MessageCode: "TRANSACTION_OTA_FAILED",
							Status:      false,
							Error:       err,
						}
						return
					}
					// end initialize new ticket

					// start creating related ticket data
					for _, item := range qrData {
						newTickdetID := uuid.NewV4()

						// start creating new tickdet
						newTickDet := entities.TickDetModel{
							Tickdet_id:      newTickdetID,
							Tickdet_tick_id: newTickID,
							Tickdet_trf_id:  item.TrfID,
							Tickdet_trftype: item.TrfType,
							Tickdet_amount:  item.TrfAmount,
							Tickdet_qty:     1,
							Tickdet_total:   item.TrfAmount,
							Tickdet_qr:      item.QR,
							Ext:             `{"void": {"status": false}, "refund": {"status": false}, "cashback": {"status": false}, "nationality": "ID"}`,
						}

						if err := tx.Create(&newTickDet).Error; err != nil {
							chResp <- ResponseDTO{
								Code:        http.StatusInternalServerError,
								Message:     err.Error(),
								MessageCode: "TRANSACTION_OTA_FAILED",
								Status:      false,
								Error:       err,
							}
							return
						}
						// end creating new tickdet

						// start fetching ticklist addition by given trf id
						var ticktlistAdditions []entities.TickListAddition
						if err := tx.Raw(`
						SELECT
							mt2.trfdet_mtick_id,
							mg.group_mid
						FROM master_tariff mt
						JOIN master_tariffdet mt2 ON mt2.trfdet_trf_id = mt.trf_id
						JOIN master_ticket mt3 ON mt3.mtick_id = mt2.trfdet_mtick_id
						JOIN master_group mg ON mg.group_id = mt3.mtick_group_id
						WHERE mt.trf_id = ?
						`, item.TrfID).Scan(&ticktlistAdditions).Error; err != nil {
							chResp <- ResponseDTO{
								Code:        http.StatusInternalServerError,
								Message:     err.Error(),
								MessageCode: "TRANSACTION_OTA_FAILED",
								Status:      false,
								Error:       err,
							}
							return
						}
						// end fetching ticklist addition by given trf id

						// start creating new ticklist
						for _, ticktlistAddition := range ticktlistAdditions {
							visitDate, err := time.Parse("2006-01-02", req.VisitDate)
							if err != nil {
								chResp <- ResponseDTO{
									Code:        http.StatusInternalServerError,
									Message:     err.Error(),
									MessageCode: "TRANSACTION_OTA_FAILED",
									Status:      false,
									Error:       err,
								}
								return
							}
							expiryDate := visitDate.Add(time.Hour * 24).Add(time.Second * -1)
							newTickList := entities.TickListModel{
								Ticklist_id:         uuid.NewV4(),
								Ticklist_tickdet_id: newTickdetID,
								Ticklist_mtick_id:   ticktlistAddition.TrfdetMtickID,
								Ticklist_expire:     expiryDate.Format("2006-01-02 15:04:05"),
								Ticklist_mid:        ticktlistAddition.GroupMid,
							}

							if err := tx.Create(&newTickList).Error; err != nil {
								chResp <- ResponseDTO{
									Code:        http.StatusInternalServerError,
									Message:     err.Error(),
									MessageCode: "TRANSACTION_OTA_FAILED",
									Status:      false,
									Error:       err,
								}
								return
							}
						}
						// end creating new ticklist
					}
					// end creating related ticket data
				}
				// end main process
			}

			if err := tx.Commit().Error; err != nil {
				chResp <- ResponseDTO{
					Code:        http.StatusInternalServerError,
					Message:     err.Error(),
					MessageCode: "TRANSACTION_OTA_FAILED",
					Status:      false,
					Error:       err,
				}
				return
			}

			chResp <- ResponseDTO{Error: nil}
		}(batch)
	}

	wg.Wait()
	close(chResp)

	for res := range chResp {
		if res.Error != nil {
			return resp, res.Code, res.Message, res.MessageCode, res.Status
		}
	}

	return resp, http.StatusOK, "Transaction success", "TRANSACTION_OTA_SUCCESS", true
}
