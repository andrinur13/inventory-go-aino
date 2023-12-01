package entities

import "time"

type OtaInventory struct {
	ID              string
	InventoryNumber string
	AgentID         int
	AgentName       string
	CreatedAt       *time.Time
}

func (OtaInventory) TableName() string {
	return "ota_inventory"
}

type OtaInventoryDetail struct {
	ID             string
	AgentID        int
	OtaInventoryID string
	GroupID        int
	GroupName      string
	TrfID          int
	TrfName        string
	ExpiryDate     *time.Time
	QR             *string
	TrfAmount      float32
	QrPrefix       *string
	RedeemDate     *time.Time
	VoidDate       *time.Time
	TrfType        string
	GroupMid       string
}

func (OtaInventoryDetail) TableName() string {
	return "ota_inventory_detail"
}

type TickListAddition struct {
	TrfdetMtickID int
	GroupMid      string
}
