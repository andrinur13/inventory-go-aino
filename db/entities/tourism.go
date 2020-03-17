package entities

import (
	uuid "github.com/satori/go.uuid"
)

type TarifCheck struct {
	Trf_name          string
	Trf_value         string
	Trf_currency_code string
	Trfftype_name     string
	Trfftype_id       int
	Trf_qty           int
	Day               int
	Hour              int
	Minute            int
	Trf_tax           float32
	Trf_assurance     float32
	Trf_fix_price     float32
	Trf_admin         float32
	Trf_others        float32
}

type GetCurrency struct {
	Curr_rate string
}

type TourismTickInsert struct {
	Tour_id             int
	Mtick_id            int
	Mtick_name          string
	Mtick_loc_start     int
	Mtick_loc_finish    int
	Mtick_merchant_code string
}

type TickModel struct {
	Tick_id             uuid.UUID
	Tick_stan           int
	Tick_number         string
	Tick_mid            string
	Tick_src_type       int
	Tick_src_id         string
	Tick_src_inv_num    string
	Tick_amount         float32
	Tick_emoney         int
	Tick_purc           string
	Tick_issuing        string
	Tick_date           string
	Tick_total_payment  float32
	Tick_payment_method string
}

func (TickModel) TableName() string {
	return "ticket"
}

type TickDetModel struct {
	Tickdet_id      uuid.UUID
	Tickdet_tick_id uuid.UUID
	Tickdet_trf_id  int
	Tickdet_trftype string
	Tickdet_amount  float32
	Tickdet_qty     int
	Tickdet_total   float32
	Tickdet_qr      string
	Ext             string
}

func (TickDetModel) TableName() string {
	return "ticketdet"
}

type TickListModel struct {
	Ticklist_id         uuid.UUID
	Ticklist_tickdet_id uuid.UUID
	Ticklist_mtick_id   int
	// Ticklist_expire     string
	Ticklist_mid string
}

func (TickListModel) TableName() string {
	return "ticketlist"
}

type TrfResponse struct {
	TickDetID      uuid.UUID `json:"ticket_id"`
	TickDetTrfID   int       `json:"ticket_trf_id"`
	TickDetTrfType string    `json:"ticket_type"`
	TickDetAmount  float32   `json:"ticket_amount"`
	TickDetQty     int       `json:"ticket_qty"`
	TickDetQr      string    `json:"ticket_qr"`
}
