package entities

type Favorite struct {
	Fav_user_id int
	Fav_data    string
	Fav_extras  string `gorm:"default:null"`
	Fav_created string
	Fav_deleted string `gorm:"default:null"`
}

func (Favorite) TableName() string {
	return "favorite"
}
