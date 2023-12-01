package entities

import "time"

type QrAggregate struct {
	TotalTicket     int
	RedeemedTicket  int
	RemainingTicket int
}

type QrItem struct {
	TrfID      int     `json:"trf_id"`
	TrfName    string  `json:"trf_name"`
	TrfAmount  float32 `json:"trf_amount"`
	QrPrefix   *string `json:"qr_prefix"`
	QR         *string `json:"qr"`
	ExpiryDate string  `json:"expiry_date"`
}

type QrDetail struct {
	InventoryNumber string
	PksNo           string
	AgentName       string
	TrfID           int
	TrfName         string
	TrfAmount       int
	QR              string
	CreatedAt       time.Time
	RedeemDate      time.Time
	ExpiryDate      time.Time
	VoidStatus      bool
	VoidDate        time.Time
	UseDate         time.Time
}
