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
