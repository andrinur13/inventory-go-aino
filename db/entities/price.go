package entities

type GetPriceReq struct {
	DestQty int     `json:"destination_qty"`
	Visit   []Visit `json:"visit"`
}

type Visit struct {
	VisitDate string     `json:"visit_date"`
	Trf       []TrfPrice `json:"tarif"`
}

type TrfPrice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

type GetPriceRes struct {
	Visit        []VisitRes `json:"visit"`
	PriceTotal   float32    `json:"price_total"`
	Discount     float32    `json:"discount"`
	PaymentTotal float32    `json:"payment_total"`
}

type VisitRes struct {
	VisitDate string        `json:"visit_date"`
	Trf       []TrfPriceRes `json:"tarif"`
}

type TrfPriceRes struct {
	TrfID              int     `json:"trf_id"`
	TrfName            string  `json:"trf_name"`
	TrfCode            string  `json:"trf_code"`
	TrfValue           float32 `json:"trf_value"`
	Currency           string  `json:"currency"`
	Qty                int     `json:"qty"`
	DiscDStatus        bool    `json:"disc_destination_status"`
	DiscDPercent       float32 `json:"disc_destination_percent"`
	DiscDestination    float32 `json:"disc_destination"`
	DiscAStatus        bool    `json:"disc_agent_status"`
	DiscAPercent       float32 `json:"disc_agent_percent"`
	DiscAgent          float32 `json:"disc_agent"`
	TotValue           float32 `json:"tot_value"`
	TotDiscDestination float32 `json:"tot_disc_destination"`
	TotDiscAgent       float32 `json:"tot_disc_agent"`
	TotPrice           float32 `json:"tot_price"`
}
