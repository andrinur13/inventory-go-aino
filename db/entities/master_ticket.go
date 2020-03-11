package entities

type MasterTicket struct {
	Trf_id            int        `json:"trf_id"`
	Trf_name          string     `json:"trf_name"`
	Trf_code          string     `json:"trf_code"`
	Trf_group_id      int        `json:"trf_group_id"`
	Trf_trftype       string     `json:"trf_trftype"`
	Trf_value         int        `json:"trf_value"`
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
