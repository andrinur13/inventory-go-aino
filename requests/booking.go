package requests

type BookingReq struct {
	Mbmid         string    `json:"mbmid"`
	BookingNumber string    `json:"booking_number"`
	BookingDate   string    `json:"booking_date"`
	PayAmt        int       `json:"pay_amount"`
	Emoney        int       `json:"emoney"`
	PayMethod     string    `json:"payment_method"`
	Trf           []bookTrf `json:"trf"`
}

type bookTrf struct {
	TrfID     int    `json:"trf_id"`
	TrfQty    int    `json:"trf_qty"`
	TrfType   string `json:"trf_trftype"`
	TrfAmount int    `json:"trf_amount"`
	TrfTotal  int    `json:"trf_total"`
}
