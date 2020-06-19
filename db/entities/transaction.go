package entities

type TripTrxModel struct {
	Tp_id           string
	Tp_number       string
	Tp_start_date   string
	Tp_end_date     string
	Tp_duration     int
	Tp_status       int
	Tp_total_amount float32
	Tp_invoice      int
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
	Tp_number       string      `json:"invoice_number"`
	Tp_invoice      int         `json:"invoice_order"`
	Tp_start_date   string      `json:"start_date"`
	Tp_end_date     string      `json:"end_date"`
	Tp_duration     int         `json:"duration"`
	Tp_status       string      `json:"status"`
	Tp_total_amount float32     `json:"total_amount"`
	Contact         *TrxContact `json:"contact"`
	Destination     string      `json:"destination"`
}

type TrxContact struct {
	Email    string `json:"email"`
	Title    string `json:"title"`
	Fullname string `json:"full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
