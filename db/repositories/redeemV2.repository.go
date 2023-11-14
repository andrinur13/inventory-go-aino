package repositories

import (
	"fmt"
	"net/http"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"

	uuid "github.com/satori/go.uuid"
)

// RedeemTicket : redeem ticket
func RedeemTicketV2(userData *entities.Users, req *requests.RedeemReqV2) (map[string]interface{}, int, string, string, bool) {
	resp := make(map[string]interface{}, 0)

	var otaInventoryDetails []entities.OtaInventoryDetail

	if err := db.DB[0].Where("qr IN (?)", req.QR).
		Or(`qr IS NULL AND qr_prefix = ANY (
			SELECT
				split_part(element, '#', 1)
			FROM unnest(
					ARRAY[?]
				) as element
		)`, req.QR).
		Preload("OtaInventory", "agent_id = ?", userData.Typeid).
		Limit(len(req.QR)).
		Find(&otaInventoryDetails).Error; err != nil {
		return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
	}

	if len(otaInventoryDetails) != len(req.QR) {
		return resp, http.StatusNotFound, "QR not found", "TRANSACTION_OTA_NOT_FOUND", false
	}

	tx := db.DB[0].Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, item := range otaInventoryDetails {
		if userData.Typeid != item.OtaInventory.AgentID {
			return resp, http.StatusBadRequest, "You are not allowed to redeem this ticket", "TRANSACTION_OTA_FAILED", false
		}

		if item.ExpiryDate.Before(time.Now()) {
			return resp, http.StatusBadRequest, "Ticket has expired", "TRANSACTION_OTA_EXPIRED", false
		}

		if item.RedeemDate != nil {
			return resp, http.StatusBadRequest, "Ticket has been redeemed", "TRANSACTION_OTA_REDEEMED", false
		}

		if item.VoidDate != nil {
			return resp, http.StatusBadRequest, "Ticket has been voided", "TRANSACTION_OTA_VOIDED", false
		}

		if item.QrPrefix == nil {
			item.QR = &req.QR[i]
		}

		now := time.Now()
		item.RedeemDate = &now

		if err := tx.Save(&item).Error; err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}

		tickStan := time.Now().UnixNano()
		microStan := tickStan / (int64(time.Millisecond) / int64(time.Nanosecond))
		tickNumber := fmt.Sprintf("TWC.5.%d.%d", item.OtaInventory.AgentID, microStan)
		newTicket := entities.TickModel{
			Tick_id:             uuid.NewV4(),
			Tick_stan:           int(tickStan),
			Tick_number:         tickNumber,
			Tick_mid:            item.GroupMid,
			Tick_src_type:       5,
			Tick_src_id:         fmt.Sprintf("%d", item.OtaInventory.AgentID),
			Tick_src_inv_num:    req.OtaOrderID,
			Tick_amount:         item.TrfAmount,
			Tick_emoney:         0,
			Tick_purc:           now.Format("2006-01-02 15:04:05"),
			Tick_issuing:        now.Format("2006-01-02 15:04:05"),
			Tick_date:           req.VisitDate,
			Tick_total_payment:  item.TrfAmount,
			Tick_payment_method: "OTA",
		}

		if err := tx.Create(&newTicket).Error; err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}

		newTickDet := entities.TickDetModel{
			Tickdet_id:      uuid.NewV4(),
			Tickdet_tick_id: newTicket.Tick_id,
			Tickdet_trf_id:  item.TrfID,
			Tickdet_trftype: item.TrfType,
			Tickdet_amount:  item.TrfAmount,
			Tickdet_qty:     1,
			Tickdet_total:   item.TrfAmount,
			Tickdet_qr:      *item.QR,
			Ext:             `{"void": {"status": false}, "refund": {"status": false}, "cashback": {"status": false}, "nationality": "ID"}`,
		}

		if err := tx.Create(&newTickDet).Error; err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}

		visitDate, err := time.Parse("2006-01-02", req.VisitDate)
		if err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}

		var ticktlistAddition entities.TickListAddition
		if err := tx.Raw(`
		SELECT
			mt2.trfdet_mtick_id,
			mg.group_mid
		FROM master_tariff mt
		JOIN master_tariffdet mt2 ON mt2.trfdet_trf_id = mt.trf_id
		JOIN master_ticket mt3 ON mt3.mtick_id = mt2.trfdet_mtick_id
		JOIN master_group mg ON mg.group_id = mt3.mtick_group_id
		WHERE mt.trf_id = ?
		`, item.TrfID).Scan(&ticktlistAddition).Error; err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}

		expiryDate := visitDate.Add(time.Hour * 24).Add(time.Second * -1)
		newTickList := entities.TickListModel{
			Ticklist_id:         uuid.NewV4(),
			Ticklist_tickdet_id: newTickDet.Tickdet_id,
			Ticklist_mtick_id:   ticktlistAddition.TrfdetMtickID,
			Ticklist_expire:     expiryDate.Format("2006-01-02 15:04:05"),
			Ticklist_mid:        ticktlistAddition.GroupMid,
		}

		if err := tx.Create(&newTickList).Error; err != nil {
			return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
		}
	}

	if err := tx.Commit().Error; err != nil {
		return resp, http.StatusInternalServerError, err.Error(), "TRANSACTION_OTA_FAILED", false
	}

	return resp, http.StatusOK, "Transaction success", "TRANSACTION_OTA_SUCCESS", true
}
