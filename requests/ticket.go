package requests

import "twc-ota-api/db/entities"

type GenerateTicket struct {
	MerchantID    string   `json:"merchant_id"`
	MerchantCode  string   `json:"merchant_code"`
	MerchantName  string   `json:"merchant_name"`
	DeviceID      string   `json:"device_id"`
	DeviceCode    string   `json:"device_code"`
	DeviceName    string   `json:"device_name"`
	DatePurchase  string   `json:"date_purchase"`
	SrcType       string   `json:"src_type"`
	PayChannel    string   `json:"pay_channel"`
	NoTransaction string   `json:"no_transaction"`
	NoRef         string   `json:"no_ref"`
	Mode          string   `json:"tarif_mode"`
	Username      string   `json:"username"`
	Trf           []Tariff `json:"trf"`
}

type Tariff struct {
	TrfID string `json:"trf_id"`
	// TrfName    string `json:"trf_name"`
	// TrfTrftype string `json:"trf_trftype"`
	TrfQty int `json:"trf_qty"`
}

type Validate struct {
	MerchantID   string `json:"merchant_id"`
	MerchantCode string `json:"merchant_code"`
	MerchantName string `json:"merchant_name"`
	DeviceID     string `json:"device_id"`
	DeviceCode   string `json:"device_code"`
	DeviceName   string `json:"device_name"`
	QrCode       string `json:"qr_code"`
	QTY          int    `json:"qty"`
	Location     int    `json:"loc_id"`
}

type ValidateMember struct {
	// MerchantCode string `json:"merchant_code"`
	// MerchantName string `json:"merchant_name"`
	DeviceID     string `json:"device_id"`
	DeviceCode   string `json:"device_code"`
	DeviceName   string `json:"device_name"`
	SerialNumber string `json:"serial_number"`
}

type Void struct {
	QrCode    string `json:"qr_code"`
	TrxNumber string `json:"trx_number"`
}

type Reversal struct {
	QrCode string `json:"qr_code"`
	LocID  int    `json:"loc_id"`
}

type ValidateIN struct {
	MerchantID   string `json:"merchant_id"`
	MerchantCode string `json:"merchant_code"`
	MerchantName string `json:"merchant_name"`
	DeviceID     string `json:"device_id"`
	DeviceCode   string `json:"device_code"`
	DeviceName   string `json:"device_name"`
	QrCode       string `json:"qr_code"`
	QTY          int    `json:"qty"`
	LocID        int    `json:"loc_id"`
	DateValidate string `json:"date_validate"`
}

type LimitRequest struct {
	MerchantCode string `json:"merchant_code"`
}

type Transaction struct {
	// MerchantCode string `json:"merchant_code"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
}

type GetQrRequest struct {
	Limit int `form:"limit"`
}

type GetQrResponse struct {
	TotalTicket     int               `json:"total_ticket"`
	RedeemedTicket  int               `json:"redeemed_ticket"`
	RemainingTicket int               `json:"remaining_ticket"`
	QrData          []entities.QrItem `json:"qr_data"`
}
