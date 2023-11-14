package requests

// RedeemReq : request redeem
type RedeemReq struct {
	BookNumber string `json:"book_number"`
}

type RedeemReqV2 struct {
	OtaOrderID string   `json:"ota_order_id" binding:"required"`
	VisitDate  string   `json:"visit_date" binding:"required"`
	QR         []string `json:"qr" binding:"required"`
}
