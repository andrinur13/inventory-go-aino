package entities

import uuid "github.com/satori/go.uuid"

type CheckOutReq struct {
	Header Header `json:"header"`
	Trip   []Trip `json:"trip"`
	Person Person `json:"person"`
}

type Header struct {
	Order       string  `json:"invoice_order"`
	StartDate   string  `json:"startdate"`
	EndDate     string  `json:"enddate"`
	Duration    int     `json:"duration"`
	InvNumber   string  `json:"inv_number"`
	TotalAmount float32 `json:"total_amount"`
	Contact     Contact `json:"contact"`
}
type Contact struct {
	Title    string `json:"title"`
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type Trip struct {
	Day         int           `json:"day"`
	Tanggal     string        `json:"tanggal"`
	Date        string        `json:"date"`
	Destination []Destination `json:"destination"`
}

type Destination struct {
	Mid          string  `json:"group_mid"`
	GroupName    string  `json:"group_name"`
	Duration     int     `json:"duration"`
	TrfAdult     float32 `json:"trf_adult"`
	TrfChild     float32 `json:"trf_child"`
	Operational  string  `json:"operational"`
	Total        float32 `json:"total"`
	Trf_id_adult int     `json:"trf_id_adult"`
	Trf_id_child int     `json:"trf_id_child"`
}

type Person struct {
	Adult []Persons `json:"adult"`
	Child []Persons `json:"child"`
}

type Persons struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	TypeID string `json:"typeid"`
}

type TripModel struct {
	Tp_id           uuid.UUID
	Tp_stan         int
	Tp_number       string
	Tp_src_type     int
	Tp_start_date   string
	Tp_end_date     string
	Tp_duration     int
	Tp_status       int
	Tp_user_id      int
	Tp_adult        int
	Tp_child        int
	Tp_contact      string
	Tp_total_amount float32
	Tp_extras       string
	Tp_invoice      int
	Created_at      string
	Tp_agent_id     int
}

func (TripModel) TableName() string {
	return "trip_planner"
}

type PersonModel struct {
	Tpp_id     uuid.UUID
	Tpp_tp_id  uuid.UUID
	Tpp_name   string
	Tpp_type   int
	Tpp_qr     string
	Tpp_extras string
	Created_at string
}

func (PersonModel) TableName() string {
	return "trip_planner_person"
}

type DestinationModel struct {
	Tpd_id        uuid.UUID
	Tpd_tpp_id    uuid.UUID
	Tpd_group_mid string
	Tpd_trf_id    int
	Tpd_amount    float32
	Tpd_date      string
	Tpd_exp_date  string
	Tpd_day       int
	Tpd_duration  int
	Tpd_extras    string
	Created_at    string
}

func (DestinationModel) TableName() string {
	return "trip_planner_destination"
}

type QRTripRes struct {
	Type   string    `json:"person_type"`
	QRcode string    `json:"qr_code"`
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Title  string    `json:"title"`
	TypeID string    `json:"typeid"`
	Trip   []TripDay `json:"trip"`
}

type TripDay struct {
	Day         int    `json:"day"`
	Tanggal     string `json:"tanggal"`
	Date        string `json:"date"`
	ExpiredDate string `json:"expired_date"`
	Destination []Dest `json:"destination"`
}

type Dest struct {
	Mid         string  `json:"group_mid"`
	GroupName   string  `json:"group_name"`
	Duration    int     `json:"duration"`
	Operational string  `json:"operational"`
	Amount      float32 `json:"amount"`
	Trf_id      int     `json:"trf_id"`
}
