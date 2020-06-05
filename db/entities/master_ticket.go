package entities

type MasterTicket struct {
	Trf_id            int        `json:"trf_id"`
	Trf_name          string     `json:"trf_name"`
	Trf_code          string     `json:"trf_code"`
	Trf_group_id      int        `json:"trf_group_id"`
	Trf_trftype       string     `json:"trf_trftype"`
	Trf_value         float32    `json:"trf_value"`
	Trf_start_date    string     `json:"trf_start_date"`
	Trf_end_date      string     `json:"trf_end_date"`
	Trf_priority      int        `json:"trf_priority"`
	Trf_release       string     `json:"trf_release"`
	Trf_currency_code string     `json:"trf_currency_code"`
	Trf_qty           int        `json:"trf_qty"`
	Day               string     `json:"-"`
	Begin_time        string     `json:"-"`
	End_time          string     `json:"-"`
	Card_type         string     `json:"-"`
	Expired_qr        int        `json:"-"`
	Trf_condition     *Condition `json:"trf_condition"`
	Ticket            []Ticket   `json:"ticket"`
}

type Ticket struct {
	Mtick_name string `json:"mtick_name"`
	Mtick_code string `json:"mtick_code"`
	Mtick_type string `json:"mtick_type"`
	Mtick_cat  string `json:"mtick_cat"`
}

type Condition struct {
	Day        []string `json:"day"`
	Begin_time string   `json:"beginTime"`
	End_time   string   `json:"endTime"`
	Card_type  string   `json:"cardType"`
	Expired_qr int      `json:"expiredQr"`
}

type GrupModel struct {
	Group_id   int
	Group_mid  string
	Group_name string
}

func (GrupModel) TableName() string {
	return "master_group"
}

type TariffModel struct {
	Trf_id    int
	Trf_name  string
	Trf_code  string
	Trf_value float32
}

func (TariffModel) TableName() string {
	return "master_tariff"
}

type TicketModel struct {
	Mtick_id       int
	Mtick_name     string
	Group          GrupModel `gorm:"foreignkey:Mtick_group_id;association_foreignkey:Group_id"`
	Mtick_group_id int
}

func (TicketModel) TableName() string {
	return "master_ticket"
}

type TariffDetModel struct {
	Tarif           TariffModel `gorm:"foreignkey:Trfdet_trf_id;association_foreignkey:Trf_id"`
	Ticket          TicketModel `gorm:"foreignkey:Trfdet_mtick_id;association_foreignkey:Mtick_id"`
	Trfdet_trf_id   int
	Trfdet_mtick_id int
}

func (TariffDetModel) TableName() string {
	return "master_tariffdet"
}
