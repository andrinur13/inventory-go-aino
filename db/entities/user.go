package entities

//Users : model for get data user
type Users struct {
	ID       int    `json:"-`
	Email    string `json:"email"`
	Password string `json:"-"`
	Typeid   int    `json:"agent_id"`
}
