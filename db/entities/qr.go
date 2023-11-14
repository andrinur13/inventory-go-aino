package entities

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
