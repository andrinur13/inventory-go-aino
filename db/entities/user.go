package entities

//Users : model for get data user
type Users struct {
	ID           int    `json:"-" gorm:"primary_key"`
	Name         string `json:"nama_depan"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Type         string `json:"type"`
	Typeid       int    `json:"agent_id"`
	Users_extras string `json:"-"`
	Created_at   string `json:"-"`
}

type UserReq struct {
	Name          string `json:"nama_depan"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Type          string `json:"type"`
	Typeid        int    `json:"agent_id"`
	NationalityID int    `json:"nationality_id"`
}
