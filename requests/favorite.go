package requests

type FavReq struct {
	Name           string    `json:"name"`
	Image          string    `json:"image"`
	ImageURL       string    `json:"image_url"`
	DestinationQty int       `json:"destination_qty"`
	Duration       int       `json:"duration"`
	Adult          int       `json:"adult"`
	Child          int       `json:"child"`
	NationalityID  int       `json:"nationality_id"`
	Bruto          float64   `json:"price_bruto"`
	Netto          float64   `json:"price_netto"`
	Disc           float64   `json:"price_disc"`
	Data           []FavData `json:"data"`
}

type FavData struct {
	Day string   `json:"day"`
	Trf []FavTrf `json:"tarif"`
}

type FavTrf struct {
	TrfID    int     `json:"trf_id"`
	TrfCode  string  `json:"trf_code"`
	TrfQty   int     `json:"trf_qty"`
	TrfNetto float32 `json:"trf_netto"`
}

type FavDelete struct {
	FavID string `json:"paket_id"`
}

type FavUploadImage struct {
	FavID string `json:"paket_id"`
	Image string `json:"image"`
}
