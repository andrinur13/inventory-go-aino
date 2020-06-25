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

type TrxRes struct {
	SourceType int       `json:"source_type"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	DestQty    int       `json:"destination_qty"`
	Trip       []TrxTrip `json:"trip"`
	Customer   []TrxCust `json:"customer"`
}

type TrxTrip struct {
	TripDate string    `json:"trip_date"`
	TripDay  int       `json:"trip_day"`
	Ticket   []TrxTick `json:"ticket"`
}

type TrxCust struct {
	Nationality string `json:"cs_nationality"`
	Region      string `json:"cs_region"`
	IDType      string `json:"cs_id_type"`
	IDNumber    string `json:"cs_id_number"`
	Name        string `json:"cs_name"`
	Type        string `json:"cs_type"`
	Title       string `json:"cs_title"`
	Email       string `json:"cs_email"`
	Phone       string `json:"cs_phone"`
	IsPic       bool   `json:"cs_is_pic"`
}

type TrxTick struct {
	Mmid         string `json:"mmid"`
	SiteDuration string `json:"site_duration"`
	TrfID        int    `json:"trf_id"`
	TrfName      string `json:"trf_name"`
	TrfQty       int    `json:"trf_qty"`
}
