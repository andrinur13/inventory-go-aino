package requests

type FavReq struct {
	Name  string    `json:"name"`
	Bruto float64   `json:"price_bruto"`
	Netto float64   `json:"price_netto"`
	Disc  float64   `json:"price_disc"`
	Data  []FavData `json:"data"`
}

type FavData struct {
	Day string   `json:"day"`
	Trf []FavTrf `json:"tarif"`
}

type FavTrf struct {
	TrfID    int     `json:"trf_id"`
	TrfNetto float32 `json:"trf_netto"`
}

type FavDelete struct {
	FavID string `json:"paket_id"`
}
