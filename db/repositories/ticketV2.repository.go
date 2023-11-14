package repositories

import (
	"net/http"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"
)

// GetQRV2 : get list ticket code
func GetQRV2(userData *entities.Users, query *requests.GetQrRequest) (*requests.GetQrResponse, int, string, string, bool) {
	var (
		aggregate          entities.QrAggregate
		otaInventoryDetail []entities.OtaInventoryDetail
		qrData             []entities.QrItem
	)

	if err := db.DB[0].Raw(`
	SELECT
		COUNT(CASE WHEN oid2.void_date IS NULL THEN 1 END) AS total_ticket,
		COUNT(CASE WHEN oid2.redeem_date IS NOT NULL THEN 1 END) AS redeemed_ticket,
		COUNT(CASE WHEN oid2.void_date IS NULL AND oid2.redeem_date is NULL THEN 1 END) AS remaining_ticket
	FROM ota_inventory_detail oid2
	JOIN ota_inventory oi ON oi.id = oid2.ota_inventory_id
	AND oi.agent_id = ?`, userData.Typeid).Scan(&aggregate).Error; err != nil {
		return nil, http.StatusInternalServerError, err.Error(), "LIST_CODE_500", false
	}

	if err := db.DB[0].Select("trf_id, trf_name, trf_amount, qr_prefix, qr, expiry_date").
		Joins("JOIN ota_inventory ON ota_inventory.id = ota_inventory_detail.ota_inventory_id").
		Where("agent_id = ? AND redeem_date IS NULL AND void_date IS NULL", userData.Typeid).Limit(query.Limit).Find(&otaInventoryDetail).Error; err != nil {
		return nil, http.StatusInternalServerError, err.Error(), "LIST_CODE_500", false
	}

	for _, item := range otaInventoryDetail {
		qrData = append(qrData, entities.QrItem{
			TrfID:      item.TrfID,
			TrfName:    item.TrfName,
			TrfAmount:  item.TrfAmount,
			QrPrefix:   item.QrPrefix,
			QR:         item.QR,
			ExpiryDate: item.ExpiryDate.Format("2006-01-02 15:04:05"),
		})
	}

	resp := &requests.GetQrResponse{
		TotalTicket:     aggregate.TotalTicket,
		RedeemedTicket:  aggregate.RedeemedTicket,
		RemainingTicket: aggregate.RemainingTicket,
		QrData:          qrData,
	}

	return resp, http.StatusOK, "Get ticket code list successfully", "LIST_CODE_200", true
}
