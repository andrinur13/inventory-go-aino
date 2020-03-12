package requests

type BookingReq struct {
	Mbmid         string    `json:"mbmid"`
	BookingNumber string    `json:"booking_number"`
	BookingDate   string    `json:"booking_date"`
	PayAmt        float32   `json:"pay_amount"`
	Emoney        int       `json:"emoney"`
	PayMethod     string    `json:"payment_method"`
	Username      string    `json:"customer_username"`
	Email         string    `json:"customer_email"`
	Phone         string    `json:"customer_phone"`
	Note          string    `json:"customer_note"`
	Trf           []bookTrf `json:"trf"`
}

type bookTrf struct {
	TrfID     int     `json:"trf_id"`
	TrfQty    int     `json:"trf_qty"`
	TrfType   string  `json:"trf_trftype"`
	TrfAmount float32 `json:"trf_amount"`
	TrfTotal  float32 `json:"trf_total"`
}
