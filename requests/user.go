package requests

type UpdatePassword struct {
	OldPwd  string `json:"old_password"`
	NewPwd  string `json:"new_password"`
	ConfPwd string `json:"confirm_password"`
}

type ResetPassword struct {
	Email string `json:"email"`
}

type UpdateResetPassword struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	NewPwd  string `json:"password"`
	ConfPwd string `json:"confirm_password"`
}
