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
	Agent_name      string
	Email           string
	Title           string
	Fullname        string
	Phone           string
	Address         string
	Created_at      string
}

func (TripTrxModel) TableName() string {
	return "trip_planner"
}

type TripPersonTrxModel struct {
	Tpp_id    string
	Tpp_name  string
	Type      string
	Tpp_qr    string
	Id_number string
	Title     string
	Type_id   string
}

func (TripPersonTrxModel) TableName() string {
	return "trip_planner_person"
}

type TripDestinationTrxModel struct {
	Trf_name      string  `json:"trf_name"`
	Tpd_group_mid string  `json:"mid"`
	Group_name    string  `json:"merchant_name"`
	Tpd_amount    float32 `json:"amount"`
	Tpd_duration  int     `json:"duration"`
	Tpd_date      string  `json:"visit_date"`
	Tpd_exp_date  string  `json:"expired_date"`
}

func (TripDestinationTrxModel) TableName() string {
	return "trip_planner_destination"
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
	Agent_name      string      `json:"agent_name"`
	Destination     string      `json:"destination"`
	Person          []TrxPerson `json:"person"`
}

type TrxContact struct {
	Email    string `json:"email"`
	Title    string `json:"title"`
	Fullname string `json:"full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type TrxPerson struct {
	Tpp_name    string                    `json:"name"`
	Type        string                    `json:"type"`
	Tpp_qr      string                    `json:"qr_code"`
	Id_number   string                    `json:"number_id"`
	Type_id     string                    `json:"type_id"`
	Title       string                    `json:"title"`
	Destination []TripDestinationTrxModel `json:"visit"`
}

type TrpTrxModel struct {
	Tp_id             uuid.UUID
	Tp_stan           int
	Tp_number         string
	Tp_src_type       int
	Tp_start_date     string
	Tp_end_date       string
	Tp_duration       int
	Tp_status         int
	Tp_user_id        int
	Tp_contact        string
	Tp_total_amount   float32
	Created_at        string
	Tp_agent_id       int
	Tp_payment_method string `gorm:"DEFAULT:null"`
	Updated_at        string `gorm:"DEFAULT:null"`
}

func (TrpTrxModel) TableName() string {
	return "trip_planner"
}

type GetExp struct {
	Expired   int
	Group_mid string
	Duration  string
}

func (GetExp) TableName() string {
	return "master_tariff"
}

type UpdateTrxResp struct {
	InvNumber string `json:"booking_number"`
	Status    string `json:"status"`
}

type DestinationTrxModel struct {
	Group_name string  `json:"destinasi"`
	Trf_name   string  `json:"tarif"`
	Tpd_amount float32 `json:"netto"`
	Bruto      float32 `json:"bruto"`
	Disc       float32 `json:"disc"`
}

func (DestinationTrxModel) TableName() string {
	return "trip_planner_destination"
}

type RespTrxNum struct {
	Tp_number       string          `json:"number"`
	Tp_trx_date     string          `json:"trxdate"`
	Tp_start_date   string          `json:"startdate"`
	Tp_end_date     string          `json:"enddate"`
	Tp_duration     int             `json:"duration"`
	Tp_status       string          `json:"status"`
	Tp_total_amount float32         `json:"amount"`
	Agent           string          `json:"agent"`
	Fullname        string          `json:"piccust"`
	Person          []RespPersonNum `json:"person"`
}

type RespPersonNum struct {
	Tpp_name    string                `json:"cust"`
	Type        string                `json:"type"`
	Tpp_qr      string                `json:"QR"`
	Destination []DestinationTrxModel `json:"detail"`
}
