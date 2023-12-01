package entities

import "time"

type QrAggregate struct {
	TotalTicket     int
	RedeemedTicket  int
	RemainingTicket int
}

type QrItem struct {
	TrfID      int
	TrfName    string
	TrfAmount  float32
	QrPrefix   *string
	QR         *string
	ExpiryDate string
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
