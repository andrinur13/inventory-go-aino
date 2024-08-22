package repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/logger"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go.elastic.co/apm/module/apmgorm/v2"
	"go.elastic.co/apm/v2"
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
func RedeemTicketV2(ctx context.Context, userData *entities.Users, req *requests.RedeemReqV2) (map[string]interface{}, int, string, string, bool) {
	resp := make(map[string]interface{}, 0)
	dbConn0 := apmgorm.WithContext(ctx, db.DB[0])

	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return resp, http.StatusBadRequest, err.Error(), "TRANSACTION_OTA_DATE_INVALID", false
	}

	today := time.Now().Truncate(24 * time.Hour)
	if visitDate.Before(today) {
		return resp, http.StatusBadRequest, "Visit date must be today or later", "TRANSACTION_OTA_DATE_INVALID", false
	}

	if len(req.QR) > 100 {
		return resp, http.StatusBadRequest, "Maximum 50 QR per request", "TRANSACTION_OTA_MAX_QR", false
	}

	sort.Strings(req.QR)
	batchSize := 10
	expectedResponses := len(req.QR) / batchSize
	if len(req.QR)%batchSize != 0 {
		expectedResponses++
	}

	wg := new(sync.WaitGroup)
	chResp := make(chan ResponseDTO, expectedResponses)
	mtx := new(sync.Mutex)

	// start main loop
	batchLoopSpan, _ := apm.StartSpan(ctx, "RedeemReqV2.batchLoopSpan", "repository")

	for start, end := 0, 0; start < len(req.QR); start = end {
		end = start + batchSize
		if end > len(req.QR) {
			end = len(req.QR)
		}

		batch := req.QR[start:end]
		wg.Add(1)
		go func(innerBatch []string) {
			defer wg.Done()

			mtx.Lock()
			defer mtx.Unlock()

			// start construct qr prefix
			qrPrefixes := constructQrPrefix(innerBatch)
			// end construct qr prefix

			// start looping qr prefix
			for _, qrPrefix := range qrPrefixes {
				// start fetching ota inventory details by given qr
				var otaInventoryDetails []entities.OtaInventoryDetail

				err := db.WithRetry("otaInventoryDetails", fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), dbConn0, func(innerDb *gorm.DB) error {
					return innerDb.Raw(`
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
						Scan(&otaInventoryDetails).Error
				})

				if err != nil {
					logger.Error("Error when fetching ota inventory details", "500", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), err)

					msg := err.Error()
					if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" || strings.Contains(err.Error(), "connection refused") {
						msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
					}

					chResp <- ResponseDTO{
						Code:        http.StatusInternalServerError,
						Message:     msg,
						MessageCode: "TRANSACTION_OTA_FAILED",
						Status:      false,
						Error:       err,
					}

					return
				}

				if len(otaInventoryDetails) == 0 {
					logger.Warning("QR not found", "404", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count))
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
					logger.Warning("QR maximum exceeded", "400", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count))
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

						logger.Warning(fmt.Sprintf("Ticket %s has expired", qr), "400", false, fmt.Sprintf("agent_id: %d, qr: %s, group_mid: %s", userData.Typeid, qr, item.GroupMid))
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

					err := db.WithTransactionRetry(dbConn0, func(tx *gorm.DB) error {
						// start updating and appending qr data by given mid
						for i, otaInventoryDetail := range otaInventoryDetails {
							if otaInventoryDetail.GroupMid == mid {
								// start checking & updating ota inventory detail
								if otaInventoryDetail.QrPrefix != "" {
									otaInventoryDetail.QR = qrPrefix.Qr[i]
								}

								if err := tx.Model(&otaInventoryDetail).Where("id = ?", otaInventoryDetail.ID).Update("redeem_date", &now).Update("qr", otaInventoryDetail.QR).Error; err != nil {
									logger.Error("Error when updating ota inventory detail", "500", false, fmt.Sprintf("%+v", otaInventoryDetail), err)
									chResp <- ResponseDTO{
										Code:        http.StatusInternalServerError,
										Message:     err.Error(),
										MessageCode: "TRANSACTION_OTA_FAILED",
										Status:      false,
										Error:       err,
									}
									return err
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

						tx.NewRecord(newTicket)

						if err := tx.Create(&newTicket).Error; err != nil {
							logger.Error("Error when creating new ticket", "500", false, fmt.Sprintf("%+v", newTicket), err)
							chResp <- ResponseDTO{
								Code:        http.StatusInternalServerError,
								Message:     err.Error(),
								MessageCode: "TRANSACTION_OTA_FAILED",
								Status:      false,
								Error:       err,
							}
							return err
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

							tx.NewRecord(newTickDet)

							if err := tx.Create(&newTickDet).Error; err != nil {
								logger.Error("Error when creating new tickdet", "500", false, fmt.Sprintf("%+v", newTickDet), err)
								chResp <- ResponseDTO{
									Code:        http.StatusInternalServerError,
									Message:     err.Error(),
									MessageCode: "TRANSACTION_OTA_FAILED",
									Status:      false,
									Error:       err,
								}
								return err
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
								logger.Error("Error when fetching ticklist addition by given trf id", "500", false, fmt.Sprintf("trf_id: %d", item.TrfID), err)
								chResp <- ResponseDTO{
									Code:        http.StatusInternalServerError,
									Message:     err.Error(),
									MessageCode: "TRANSACTION_OTA_FAILED",
									Status:      false,
									Error:       err,
								}
								return err
							}
							// end fetching ticklist addition by given trf id

							// start creating new ticklist
							for _, ticktlistAddition := range ticktlistAdditions {
								visitDate, err := time.Parse("2006-01-02", req.VisitDate)
								if err != nil {
									logger.Error("Error when parsing visit date", "500", false, fmt.Sprintf("visit_date: %s", req.VisitDate), err)
									chResp <- ResponseDTO{
										Code:        http.StatusInternalServerError,
										Message:     err.Error(),
										MessageCode: "TRANSACTION_OTA_FAILED",
										Status:      false,
										Error:       err,
									}
									return err
								}
								expiryDate := visitDate.Add(time.Hour * 24).Add(time.Second * -1)
								newTickList := entities.TickListModel{
									Ticklist_id:         uuid.NewV4(),
									Ticklist_tickdet_id: newTickdetID,
									Ticklist_mtick_id:   ticktlistAddition.TrfdetMtickID,
									Ticklist_expire:     expiryDate.Format("2006-01-02 15:04:05"),
									Ticklist_mid:        ticktlistAddition.GroupMid,
								}

								tx.NewRecord(newTickList)

								if err := tx.Create(&newTickList).Error; err != nil {
									logger.Error("Error when creating new ticklist", "500", false, fmt.Sprintf("%+v", newTickList), err)
									chResp <- ResponseDTO{
										Code:        http.StatusInternalServerError,
										Message:     err.Error(),
										MessageCode: "TRANSACTION_OTA_FAILED",
										Status:      false,
										Error:       err,
									}
									return err
								}
							}
							// end creating new ticklist
						}
						// end creating related ticket data

						return nil
					})

					if err != nil {
						logger.Error(fmt.Sprintf("Error occured when: %s", err.Error()), "500", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), err)

						msg := err.Error()
						if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" || strings.Contains(err.Error(), "connection refused") {
							msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
						}

						chResp <- ResponseDTO{
							Code:        http.StatusInternalServerError,
							Message:     msg,
							MessageCode: "TRANSACTION_OTA_FAILED",
							Status:      false,
							Error:       err,
						}

						return
					}
				}
				// end main process
			}
			// end looping qr prefix

			chResp <- ResponseDTO{
				Code:    http.StatusOK,
				Message: "Batch processed successfully",
				Status:  true,
			}
		}(batch)
	}

	batchLoopSpan.End()
	// end main loop

	wg.Wait()
	close(chResp)

	for res := range chResp {
		if res.Error != nil {
			return resp, res.Code, res.Message, res.MessageCode, res.Status
		} else {
			return resp, http.StatusOK, "Transaction success", "TRANSACTION_OTA_SUCCESS", true
		}
	}

	return resp, http.StatusOK, "Transaction success", "TRANSACTION_OTA_SUCCESS", true
}
