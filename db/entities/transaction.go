package entities

import uuid "github.com/satori/go.uuid"

type TripTrxModel struct {
	Tp_id           string
	Tp_number       string
	Tp_start_date   string
	Tp_end_date     string
	Tp_duration     int
	Tp_status       int
	Tp_total_amount float32
	Tp_invoice      int
	Tp_agent_id     int
	Email           string
	Title           string
	Fullname        string
	Phone           string
	Address         string
}

func (TripTrxModel) TableName() string {
	return "trip_planner"
}

type TripGrupName struct {
	Group_name string
}

func (TripGrupName) TableName() string {
	return "trip_planner"
}

type TrxList struct {
	Tp_id           string      `json:"tp_id"`
	Tp_number       string      `json:"invoice_number"`
	Tp_invoice      int         `json:"invoice_order"`
	Tp_start_date   string      `json:"start_date"`
	Tp_end_date     string      `json:"end_date"`
	Tp_duration     int         `json:"duration"`
	Tp_status       int         `json:"tp_status"`
	Status_name     string      `json:"status_name"`
	Tp_total_amount float32     `json:"total_amount"`
	Contact         *TrxContact `json:"contact"`
	Tp_agent_id     int         `json:"agent_id"`
	Destination     string      `json:"destination"`
}

type TrxContact struct {
	Email    string `json:"email"`
	Title    string `json:"title"`
	Fullname string `json:"full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type TrpTrxModel struct {
	Tp_id           uuid.UUID
	Tp_stan         int
	Tp_number       string
	Tp_src_type     int
	Tp_start_date   string
	Tp_end_date     string
	Tp_duration     int
	Tp_status       int
	Tp_user_id      int
	Tp_contact      string
	Tp_total_amount float32
	Created_at      string
	Tp_agent_id     int
}

func (TrpTrxModel) TableName() string {
	return "trip_planner"
}

type GetExp struct {
	Expired int
}

func (GetExp) TableName() string {
	return "master_tariff"
}
