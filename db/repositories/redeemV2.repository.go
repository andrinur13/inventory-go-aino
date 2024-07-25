package repositories

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/logger"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	uuid "github.com/satori/go.uuid"
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
	maxRetry := 10

	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return resp, http.StatusBadRequest, err.Error(), "TRANSACTION_OTA_DATE_INVALID", false
	}

	today := time.Now().Truncate(24 * time.Hour)
	if visitDate.Before(today) {
		return resp, http.StatusBadRequest, "Visit date must be today or later", "TRANSACTION_OTA_DATE_INVALID", false
	}

	sort.Strings(req.QR)
	batchSize := 100

	batchLoopSpanTime := time.Now()
	batchLoopSpan, batchLoopCtx := apm.StartSpan(ctx, "RedeemReqV2.batchLoopSpan", "repository")

	for start, end := 0, 0; start < len(req.QR); start = end {
		end = start + batchSize
		if end > len(req.QR) {
			end = len(req.QR)
		}

		batch := req.QR[start:end]

		// checking if qr is not redeemed
		checkOtaInventoryDetailsSpan, _ := apm.StartSpan(batchLoopCtx, "checkOtaInventoryDetailsSpan", "repository")
		checkOtaInventoryDetailsSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
		checkOtaInventoryDetailsSpan.Context.SetLabel("qr", strings.Join(batch, ","))

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
			logger.Error("Error when fetching ota inventory details", "500", false, fmt.Sprintf("agent_id: %d, qr: %+v", userData.Typeid, []string{strings.Join(batch, ",")}), err)
			if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
				for i := 0; i < maxRetry; i++ {
					err = db.DB[0].Raw(`
						SELECT oid2.*
						FROM ota_inventory_detail oid2
						JOIN ota_inventory oi ON oi.id = oid2.ota_inventory_id
						WHERE oi.agent_id = ?
						AND oid2.qr IN (?)
						AND oid2.redeem_date IS NOT NULL
						LIMIT ?`, userData.Typeid, batch, len(batch)).
						Scan(&checkOtaInventoryDetails).Error
					if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
						logger.Error(fmt.Sprintf("Hitback fetching ota inventory details: %d...", i+1), "500", false, fmt.Sprintf("agent_id: %d, qr: %+v", userData.Typeid, []string{strings.Join(batch, ",")}), err)
						time.Sleep(3 * time.Second)
						continue
					} else {
						logger.Info(fmt.Sprintf("Hitback fetching ota inventory details successfully: %d...", i+1), "200", true, "", fmt.Sprintf("agent_id: %d, qr: %+v", userData.Typeid, []string{strings.Join(batch, ",")}))
						break
					}
				}

				if err != nil {
					logger.Error("Error when hitback fetching ota inventory details", "500", false, fmt.Sprintf("agent_id: %d, qr: %+v", userData.Typeid, []string{strings.Join(batch, ",")}), err)

					msg := err.Error()
					if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
						msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
					}

					return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
				}
			} else {
				return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
			}
		}

		if len(checkOtaInventoryDetails) > 0 {
			logger.Warning(fmt.Sprintf("QR %s has been redeemed", checkOtaInventoryDetails[0].QR), "400", false, fmt.Sprintf("agent_id: %d, qr: %+v", userData.Typeid, []string{strings.Join(batch, ",")}))
			return resp, http.StatusBadRequest, fmt.Sprintf("QR %s has been redeemed", checkOtaInventoryDetails[0].QR), "TRANSACTION_OTA_REDEEMED", false
		}
		checkOtaInventoryDetailsSpan.End()
		// end checking if qr is not redeemed

		// start construct qr prefix
		qrPrefixesSpan, _ := apm.StartSpan(ctx, "qrPrefixesSpan", "repository")
		qrPrefixesSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
		qrPrefixesSpan.Context.SetLabel("qr", strings.Join(batch, ","))

		qrPrefixes := constructQrPrefix(batch)

		qrPrefixesSpan.End()
		// end construct qr prefix

		// start looping qr prefix
		qrPrefixLoopingSpan, qrPrefixLoopingCtx := apm.StartSpan(ctx, "qrPrefixLoopingSpan", "repository")
		qrPrefixLoopingSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
		qrPrefixLoopingSpan.Context.SetLabel("qr_prefixes", fmt.Sprintf("%+v", qrPrefixes))

		for _, qrPrefix := range qrPrefixes {
			// start fetching ota inventory details by given qr
			otaInventoryDetailsSpan, _ := apm.StartSpan(qrPrefixLoopingCtx, "otaInventoryDetailsSpan", "repository")
			otaInventoryDetailsSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
			otaInventoryDetailsSpan.Context.SetLabel("qr", qrPrefix.Qr)
			otaInventoryDetailsSpan.Context.SetLabel("qr_prefix", qrPrefix.QrPrefix)

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
				logger.Error("Error when fetching ota inventory details", "500", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), err)
				if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
					for i := 0; i < maxRetry; i++ {
						err = db.DB[0].Raw(`
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
						if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
							logger.Error(fmt.Sprintf("Hitback fetching ota inventory details: %d...", i+1), "500", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), err)
							time.Sleep(3 * time.Second)
							continue
						} else {
							logger.Info(fmt.Sprintf("Hitback fetching ota inventory details successfully: %d...", i+1), "200", true, "", fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count))
							break
						}
					}

					if err != nil {
						logger.Error("Error when hitback fetching ota inventory details", "500", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count), err)

						msg := err.Error()
						if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
							msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
						}

						return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
					}
				} else {
					return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
				}
			}

			otaInventoryDetailsSpan.Context.SetLabel("ota_inventory_details_count", len(otaInventoryDetails))
			otaInventoryDetailsSpan.Context.SetLabel("qr_prefix_count", qrPrefix.Count)

			if len(otaInventoryDetails) == 0 {
				logger.Warning("QR not found", "404", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count))
				return resp, http.StatusNotFound, "QR not found", "TRANSACTION_OTA_NOT_FOUND", false
			}

			if len(otaInventoryDetails) < qrPrefix.Count {
				logger.Warning("QR maximum exceeded", "400", false, fmt.Sprintf("agent_id: %d, qr: %s, qr_prefix: %s, limit: %d", userData.Typeid, qrPrefix.Qr, qrPrefix.QrPrefix, qrPrefix.Count))
				return resp, http.StatusBadRequest, "QR maximum exceeded", "TRANSACTION_OTA_MAX_EXCEEDED", false
			}
			otaInventoryDetailsSpan.End()
			// end fetching ota inventory details by given qr

			// start filtering mid
			midFilteringSpan, midFilteringCtx := apm.StartSpan(qrPrefixLoopingCtx, "midFilteringSpan", "repository")
			midFilteringSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
			midFilteringSpan.Context.SetLabel("ota_inventory_details", fmt.Sprintf("%+v", otaInventoryDetails))

			mids := make([]string, 0)
			for _, item := range otaInventoryDetails {
				checkExpiryDateSpan, _ := apm.StartSpan(midFilteringCtx, "checkExpiryDateSpan", "repository")
				checkExpiryDateSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
				checkExpiryDateSpan.Context.SetLabel("qr", item.QR)
				checkExpiryDateSpan.Context.SetLabel("group_mid", item.GroupMid)

				if item.ExpiryDate.Before(time.Now()) {
					qr := item.QR
					if qr == "" {
						qr = item.QrPrefix
					}

					logger.Warning(fmt.Sprintf("Ticket %s has expired", qr), "400", false, fmt.Sprintf("agent_id: %d, qr: %s, group_mid: %s", userData.Typeid, qr, item.GroupMid))
					return resp, http.StatusBadRequest, fmt.Sprintf("Ticket %s has expired", qr), "TRANSACTION_OTA_EXPIRED", false
				}

				mids = append(mids, item.GroupMid)
				checkExpiryDateSpan.End()
			}

			filteredMids := helper.RemoveDuplicateStr(mids)

			midFilteringSpan.Context.SetLabel("filtered_mids", fmt.Sprintf("%+v", filteredMids))
			midFilteringSpan.End()
			// end filtering mid

			// start main process
			mainProcessSpan, mainProcessCtx := apm.StartSpan(qrPrefixLoopingCtx, "mainProcessSpan", "repository")
			mainProcessSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
			mainProcessSpan.Context.SetLabel("filtered_mids", fmt.Sprintf("%+v", filteredMids))

			for _, mid := range filteredMids {
				// start initializing variables
				now := time.Now()
				newTickID := uuid.NewV4()
				tickAmount := new(float32)
				qrData := make([]entities.OtaInventoryDetail, 0)
				// end initializing variables

				// start updating and appending qr data by given mid
				updateQrDataSpan, updateQrDataCtx := apm.StartSpan(mainProcessCtx, "updateQrDataSpan", "repository")
				updateQrDataSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
				updateQrDataSpan.Context.SetLabel("mid", mid)
				updateQrDataSpan.Context.SetLabel("ota_inventory_details", fmt.Sprintf("%+v", otaInventoryDetails))

				for i, otaInventoryDetail := range otaInventoryDetails {
					if otaInventoryDetail.GroupMid == mid {
						// start checking & updating ota inventory detail
						checkAndUpdateOtaInventoryDetailSpan, _ := apm.StartSpan(updateQrDataCtx, "checkAndUpdateOtaInventoryDetailSpan", "repository")
						checkAndUpdateOtaInventoryDetailSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
						checkAndUpdateOtaInventoryDetailSpan.Context.SetLabel("qr", otaInventoryDetail.QR)
						checkAndUpdateOtaInventoryDetailSpan.Context.SetLabel("mid", mid)

						if otaInventoryDetail.QrPrefix != "" {
							otaInventoryDetail.QR = qrPrefix.Qr[i]
						}

						otaInventoryDetail.RedeemDate = &now

						checkAndUpdateOtaInventoryDetailSpan.Context.SetLabel("ota_inventory_detail", fmt.Sprintf("%+v", otaInventoryDetail))

						if err := db.DB[0].Save(&otaInventoryDetail).Error; err != nil {
							logger.Error("Error when updating ota inventory detail", "500", false, fmt.Sprintf("%+v", otaInventoryDetail), err)
							if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
								for i := 0; i < maxRetry; i++ {
									err = db.DB[0].Save(&otaInventoryDetail).Error
									if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
										logger.Error(fmt.Sprintf("Hitback updating ota inventory detail: %d...", i+1), "500", false, fmt.Sprintf("%+v", otaInventoryDetail), err)
										time.Sleep(3 * time.Second)
										continue
									} else {
										logger.Info(fmt.Sprintf("Hitback updating ota inventory detail successfully: %d...", i+1), "200", true, "", fmt.Sprintf("%+v", otaInventoryDetail))
										break
									}
								}

								if err != nil {
									logger.Error("Error when hitback updating ota inventory detail", "500", false, fmt.Sprintf("%+v", otaInventoryDetail), err)

									msg := err.Error()
									if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
										msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
									}

									return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
								}
							} else {
								return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
							}
						}

						qrData = append(qrData, otaInventoryDetail)
						*tickAmount += otaInventoryDetail.TrfAmount
						checkAndUpdateOtaInventoryDetailSpan.End()
						// end checking & updating ota inventory detail
					}
				}
				updateQrDataSpan.End()
				// end updating and appending qr data by given mid

				// start initialize new ticket
				insertNewTicketSpan, _ := apm.StartSpan(mainProcessCtx, "insertNewTicketSpan", "repository")
				insertNewTicketSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
				insertNewTicketSpan.Context.SetLabel("mid", mid)

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

				insertNewTicketSpan.Context.SetLabel("new_ticket", fmt.Sprintf("%+v", newTicket))

				db.DB[0].NewRecord(newTicket)

				if err := db.DB[0].Create(&newTicket).Error; err != nil {
					logger.Error("Error when creating new ticket", "500", false, fmt.Sprintf("%+v", newTicket), err)
					if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
						for i := 0; i < maxRetry; i++ {
							err = db.DB[0].Create(&newTicket).Error
							if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
								logger.Error(fmt.Sprintf("Hitback creating new ticket: %d...", i+1), "500", false, fmt.Sprintf("%+v", newTicket), err)
								time.Sleep(3 * time.Second)
								continue
							} else {
								logger.Info(fmt.Sprintf("Hitback creating new ticket successfully: %d...", i+1), "200", true, "", fmt.Sprintf("%+v", newTicket))
								break
							}
						}

						if err != nil {
							logger.Error("Error when hitback creating new ticket", "500", false, fmt.Sprintf("%+v", newTicket), err)

							msg := err.Error()
							if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
								msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
							}

							return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
						}
					} else {
						return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
					}
				}
				insertNewTicketSpan.End()
				// end initialize new ticket

				// start creating related ticket data
				createRelatedTicketDataLoopSpan, createRelatedTicketDataLoopCtx := apm.StartSpan(mainProcessCtx, "createRelatedTicketDataLoopSpan", "repository")
				createRelatedTicketDataLoopSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
				createRelatedTicketDataLoopSpan.Context.SetLabel("mid", mid)
				createRelatedTicketDataLoopSpan.Context.SetLabel("qr_data", fmt.Sprintf("%+v", qrData))

				for _, item := range qrData {
					newTickdetID := uuid.NewV4()

					// start creating new tickdet
					createNewTicdetSpan, _ := apm.StartSpan(createRelatedTicketDataLoopCtx, "createNewTicdetSpan", "repository")
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
					createNewTicdetSpan.Context.SetLabel("new_tickdet", fmt.Sprintf("%+v", newTickDet))

					db.DB[0].NewRecord(newTickDet)

					if err := db.DB[0].Create(&newTickDet).Error; err != nil {
						logger.Error("Error when creating new tickdet", "500", false, fmt.Sprintf("%+v", newTickDet), err)
						if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
							for i := 0; i < maxRetry; i++ {
								err = db.DB[0].Create(&newTickDet).Error
								if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
									logger.Error(fmt.Sprintf("Hitback creating new tickdet: %d...", i+1), "500", false, fmt.Sprintf("%+v", newTickDet), err)
									time.Sleep(3 * time.Second)
									continue
								} else {
									logger.Info(fmt.Sprintf("Hitback creating new tickdet successfully: %d...", i+1), "200", true, "", fmt.Sprintf("%+v", newTickDet))
									break
								}
							}

							if err != nil {
								logger.Error("Error when hitback creating new tickdet", "500", false, fmt.Sprintf("%+v", newTickDet), err)

								msg := err.Error()
								if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
									msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
								}

								return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
							}
						} else {
							return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
						}
					}
					createNewTicdetSpan.End()
					// end creating new tickdet

					// start fetching ticklist addition by given trf id
					fetchTicklistAdditionSpan, _ := apm.StartSpan(createRelatedTicketDataLoopCtx, "fetchTicklistAdditionSpan", "repository")
					fetchTicklistAdditionSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
					fetchTicklistAdditionSpan.Context.SetLabel("trf_id", item.TrfID)

					var ticktlistAdditions []entities.TickListAddition
					if err := db.DB[0].Raw(`
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
						if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
							for i := 0; i < maxRetry; i++ {
								err = db.DB[0].Raw(`
								SELECT
									mt2.trfdet_mtick_id,
									mg.group_mid
								FROM master_tariff mt
								JOIN master_tariffdet mt2 ON mt2.trfdet_trf_id = mt.trf_id
								JOIN master_ticket mt3 ON mt3.mtick_id = mt2.trfdet_mtick_id
								JOIN master_group mg ON mg.group_id = mt3.mtick_group_id
								WHERE mt.trf_id = ?
								`, item.TrfID).Scan(&ticktlistAdditions).Error
								if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
									logger.Error(fmt.Sprintf("Hitback fetching ticklist addition by given trf id: %d...", i+1), "500", false, fmt.Sprintf("trf_id: %d", item.TrfID), err)
									time.Sleep(3 * time.Second)
									continue
								} else {
									logger.Info(fmt.Sprintf("Hitback fetching ticklist addition by given trf id successfully: %d...", i+1), "200", true, "", fmt.Sprintf("trf_id: %d", item.TrfID))
									break
								}
							}

							if err != nil {
								logger.Error("Error when hitback fetching ticklist addition by given trf id", "500", false, fmt.Sprintf("trf_id: %d", item.TrfID), err)

								msg := err.Error()
								if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
									msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
								}

								return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
							}
						} else {
							return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
						}
					}
					fetchTicklistAdditionSpan.Context.SetLabel("ticklist_additions", fmt.Sprintf("%+v", ticktlistAdditions))
					fetchTicklistAdditionSpan.End()
					// end fetching ticklist addition by given trf id

					// start creating new ticklist
					createNewTicklistLoopSpan, createNewTicklistLoopCtx := apm.StartSpan(createRelatedTicketDataLoopCtx, "createNewTicklistLoopSpan", "repository")
					createNewTicklistLoopSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
					createNewTicklistLoopSpan.Context.SetLabel("ticklist_additions", fmt.Sprintf("%+v", ticktlistAdditions))

					for _, ticktlistAddition := range ticktlistAdditions {
						createTicklistAdditionSpan, _ := apm.StartSpan(createNewTicklistLoopCtx, "createTicklistAdditionSpan", "repository")
						createTicklistAdditionSpan.Context.SetLabel("agent_id", fmt.Sprintf("%d", userData.Typeid))
						createTicklistAdditionSpan.Context.SetLabel("trfdet_mtick_id", ticktlistAddition.TrfdetMtickID)
						createTicklistAdditionSpan.Context.SetLabel("group_mid", ticktlistAddition.GroupMid)

						visitDate, err := time.Parse("2006-01-02", req.VisitDate)
						if err != nil {
							logger.Error("Error when parsing visit date", "500", false, fmt.Sprintf("visit_date: %s", req.VisitDate), err)
							return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
						}
						expiryDate := visitDate.Add(time.Hour * 24).Add(time.Second * -1)
						newTickList := entities.TickListModel{
							Ticklist_id:         uuid.NewV4(),
							Ticklist_tickdet_id: newTickdetID,
							Ticklist_mtick_id:   ticktlistAddition.TrfdetMtickID,
							Ticklist_expire:     expiryDate.Format("2006-01-02 15:04:05"),
							Ticklist_mid:        ticktlistAddition.GroupMid,
						}

						createTicklistAdditionSpan.Context.SetLabel("new_ticklist", fmt.Sprintf("%+v", newTickList))

						db.DB[0].NewRecord(newTickList)

						if err := db.DB[0].Create(&newTickList).Error; err != nil {
							logger.Error("Error when creating new ticklist", "500", false, fmt.Sprintf("%+v", newTickList), err)
							if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
								for i := 0; i < maxRetry; i++ {
									err = db.DB[0].Create(&newTickList).Error
									if err != nil && (reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection") {
										logger.Error(fmt.Sprintf("Hitback creating new ticklist: %d...", i+1), "500", false, fmt.Sprintf("%+v", newTickList), err)
										time.Sleep(3 * time.Second)
										continue
									} else {
										logger.Info(fmt.Sprintf("Hitback creating new ticklist successfully: %d...", i+1), "200", true, "", fmt.Sprintf("%+v", newTickList))
										break
									}
								}

								if err != nil {
									logger.Error("Error when hitback creating new ticklist", "500", false, fmt.Sprintf("%+v", newTickList), err)

									msg := err.Error()
									if reflect.TypeOf(err).String() == "*net.OpError" || err.Error() == "driver: bad connection" {
										msg = "Your database connection was broken. Please contact your administrator to fix this problem. Thank you."
									}

									return resp, http.StatusInternalServerError, msg, "TRANSACTION_OTA_FAILED", false
								}
							} else {
								return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
							}
						}
						createTicklistAdditionSpan.End()
					}
					createNewTicklistLoopSpan.End()
					// end creating new ticklist
				}
				createRelatedTicketDataLoopSpan.End()
				// end creating related ticket data
			}
			mainProcessSpan.End()
			// end main process
		}
		qrPrefixLoopingSpan.End()
		// end looping qr prefix
	}
	batchLoopSpan.End()
	batchLoopSpanDuration := time.Since(batchLoopSpanTime)
	logger.Info("Batch loop span duration", "200", true, "", fmt.Sprintf("duration: %s", batchLoopSpanDuration))

	return resp, http.StatusOK, "Transaction success", "TRANSACTION_OTA_SUCCESS", true
}
