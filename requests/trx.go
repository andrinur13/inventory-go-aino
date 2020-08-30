package requests

type OfflineData struct {
	MerchantID   string      `json:"merchant_id"`
	MerchantCode string      `json:"merchant_code"`
	MerchantName string      `json:"merchant_name"`
	DeviceID     string      `json:"device_id"`
	DeviceCode   string      `json:"device_code"`
	DeviceName   string      `json:"device_name"`
	Mode         string      `json:"mode"`
	DeviceType   string      `json:"device_type"`
	Location     int         `json:"loc_id"`
	Type         string      `json:"type"`
	TrxData      []OfflineQR `json:"data"`
}

type OfflineQR struct {
	QRCode   string `json:"qr_code"`
	NoRef    string `json:"no_ref"`
	DateUsed string `json:"date_used"`
}

type OfflineMember struct {
	DeviceID   string               `json:"device_id"`
	DeviceCode string               `json:"device_code"`
	DeviceName string               `json:"device_name"`
	DeviceType string               `json:"device_type"`
	TrxData    []ValidateMemberData `json:"data"`
}

type ValidateMemberData struct {
	SerialNumber string `json:"serial_number"`
	NoRef        string `json:"no_ref"`
	DateValidate string `json:"date_validate"`
}

type TrxReq struct {
	// TotalAmount float32   `json:"total_amount"`
	SourceType int       `json:"source_type"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	DestQty    int       `json:"destination_qty"`
	Customer   []TrxCust `json:"customer"`
}

type TrxReqUpdate struct {
	Status int      `json:"status"`
	Trx    []TrxInv `json:"transaction"`
}

type TrxInv struct {
	BookingNumber string `json:"booking_number"`
	PaymentMethod string `json:"payment_method"`
}

type TrxTrip struct {
	TripDate string    `json:"trip_date"`
	TripDay  int       `json:"trip_day"`
	Ticket   []TrxTick `json:"ticket"`
}

type TrxCust struct {
	Nationality string    `json:"cs_nationality"`
	Region      string    `json:"cs_region"`
	IDType      string    `json:"cs_id_type"`
	IDNumber    string    `json:"cs_id_number"`
	Name        string    `json:"cs_name"`
	Type        string    `json:"cs_type"`
	Title       string    `json:"cs_title"`
	Email       string    `json:"cs_email"`
	Phone       string    `json:"cs_phone"`
	IsPic       bool      `json:"cs_is_pic"`
	Trip        []TrxTrip `json:"cs_trip"`
}

type TrxTick struct {
	Mmid         string  `json:"mmid"`
	SiteDuration int     `json:"site_duration"`
	TrfID        int     `json:"trf_id"`
	TrfName      string  `json:"trf_name"`
	TrfQty       int     `json:"trf_qty"`
	Amount       float32 `json:"amount"`
	Discount     float32 `json:"discount"`
	NettAmount   float32 `json:"nett_amount"`
}

type TrxResp struct {
	BookingNumber string    `json:"booking_number"`
	PayTotal      float32   `json:"payment_total"`
	Name          string    `json:"cp_name"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Cust          []TrxCust `json:"customer"`
}

type TrxVisit struct {
	TripDate string    `json:"booking_date"`
	TripDay  int       `json:"trip_day"`
	Ticket   []TrxTick `json:"ticket"`
}

type TrxQReq struct {
	Inv string `json:"invoice_number"`
}
