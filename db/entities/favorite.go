package entities

import uuid "github.com/satori/go.uuid"

type Favorite struct {
	Fav_id      uuid.UUID
	Fav_user_id int
	Fav_data    string
	Fav_extras  string `gorm:"default:null"`
	Fav_created string
	Fav_deleted string `gorm:"default:null"`
}

func (Favorite) TableName() string {
	return "favorite"
}

type FavResp struct {
	PaketID       uuid.UUID `json:"paket_id"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image_url"`
	Duration      int       `json:"duration"`
	Adult         int       `json:"adult"`
	Child         int       `json:"child"`
	NationalityID int       `json:"nationality_id"`
	Bruto         float64   `json:"price_bruto"`
	Netto         float64   `json:"price_netto"`
	Disc          float64   `json:"price_disc"`
	Data          []FavData `json:"data"`
}

type FavData struct {
	Day string   `json:"day"`
	Trf []FavTrf `json:"tarif"`
}

type FavTrf struct {
	TrfID    int     `json:"trf_id"`
	TrfName  string  `json:"trf_name"`
	TrfCode  string  `json:"trf_code"`
	TrfQty   int     `json:"trf_qty"`
	TrfNetto float32 `json:"trf_netto"`
}
