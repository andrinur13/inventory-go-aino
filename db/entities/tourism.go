package entities

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
