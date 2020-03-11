package entities

//Users : model for get data user
type Users struct {
	Email    string `json:"email"`
	Password string `json:"trf_name"`
	Typeid   int    `json:"agent_id"`
}
