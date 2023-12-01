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

// GetQRStatusV2 : get qr status
func GetQRStatusV2(userData *entities.Users, qrCode string) (*requests.GetQrStatusResponse, int, string, string, bool) {
	var qrDetail entities.QrDetail

	if err := db.DB[0].Raw(`
	SELECT
		oi.inventory_number,
		oi.pks_no,
		oi.agent_name,
		oid2.trf_id,
		oid2.trf_name,
		oid2.trf_amount,
		oid2.qr,
		oi.created_at,
		oid2.redeem_date,
		oid2.expiry_date,
		case when oid2.void_date is not null then true else false end as void_status,
		oid2.void_date
	FROM ota_inventory_detail oid2
	JOIN ota_inventory oi ON oi.id = oid2.ota_inventory_id
	AND oi.agent_id = ?
	AND oid2.qr = ?`, userData.Typeid, qrCode).Scan(&qrDetail).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, http.StatusNotFound, "QR Code not found", "STATUS_QR_404", false
		} else {
			return nil, http.StatusInternalServerError, err.Error(), "STATUS_QR_500", false
		}
	}

	qrData := requests.QrData{
		TrfID:      qrDetail.TrfID,
		TrfName:    qrDetail.TrfName,
		TrfAmount:  qrDetail.TrfAmount,
		QR:         qrDetail.QR,
		VoidStatus: qrDetail.VoidStatus,
	}

	if !qrDetail.CreatedAt.IsZero() {
		qrData.CreatedAt = qrDetail.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if !qrDetail.RedeemDate.IsZero() {
		qrData.RedeemDate = qrDetail.RedeemDate.Format("2006-01-02 15:04:05")
	}

	if !qrDetail.ExpiryDate.IsZero() {
		qrData.ExpiryDate = qrDetail.ExpiryDate.Format("2006-01-02 15:04:05")
	}

	if !qrDetail.VoidDate.IsZero() {
		qrData.VoidDate = qrDetail.VoidDate.Format("2006-01-02 15:04:05")
	}

	resp := &requests.GetQrStatusResponse{
		InventoryNumber: qrDetail.InventoryNumber,
		PksNo:           qrDetail.PksNo,
		AgentName:       qrDetail.AgentName,
		QrData:          qrData,
	}

	return resp, http.StatusOK, "Get QR status successfully", "STATUS_QR_200", true
}

// GetQRSummaryV2 : get qr summary
func GetQRSummaryV2(userData *entities.Users) (*requests.GetQrSummaryResponse, int, string, string, bool) {
	var ticketData []requests.TicketData

	if err := db.DB[0].Raw(`
	SELECT
		x.trf_id,
		x.trf_name,
		x.trf_amount,
		SUM(x.total_ticket) AS total_ticket,
		SUM(x.redeemed_ticket) AS redeemed_ticket,
		SUM(x.used_ticket) AS used_ticket,
		SUM(x.total_ticket) - SUM(x.redeemed_ticket) AS total_remaining_ticket
	FROM
		(
		SELECT
				oid.trf_id,
				oid.trf_name,
				oid.trf_amount,
				COUNT(oid.id) AS total_ticket,
				COUNT(td.tickdet_id) AS redeemed_ticket,
				CASE
					WHEN tl.ticklist_use_date IS NOT NULL THEN COUNT(td.tickdet_id) ELSE 0
				END AS used_ticket
			FROM ota_inventory oi
			LEFT JOIN ota_inventory_detail oid ON oid.ota_inventory_id = oi.id
			LEFT JOIN ticketdet td ON td.tickdet_qr = oid.qr
			LEFT JOIN ticketlist tl ON tl.ticklist_tickdet_id = td.tickdet_id
			WHERE oi.agent_id = ?
			GROUP BY
			oid.trf_id,
			oid.trf_name,
			oid.trf_amount,
			tl.ticklist_use_date
		) x
	GROUP BY
		x.trf_id,
		x.trf_name,
		x.trf_amount
	ORDER BY x.trf_name`, userData.Typeid).Scan(&ticketData).Error; err != nil {
		return nil, http.StatusInternalServerError, err.Error(), "LIST_CODE_500", false
	}

	var (
		totalTicket     int
		redeemedTicket  int
		remainingTicket int
		usedTicket      int
	)
	for _, item := range ticketData {
		totalTicket += item.TotalTicket
		redeemedTicket += item.RedeemedTicket
		remainingTicket += item.RemainingTicket
		usedTicket += item.UsedTicket
	}

	resp := &requests.GetQrSummaryResponse{
		TotalTicket:     totalTicket,
		RedeemedTicket:  redeemedTicket,
		RemainingTicket: remainingTicket,
		UsedTicket:      usedTicket,
		TicketData:      ticketData,
	}

	return resp, http.StatusOK, "Get ticket code list successfully", "LIST_CODE_200", true
}
