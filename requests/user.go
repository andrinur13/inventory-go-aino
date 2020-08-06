package requests

type UpdatePassword struct {
	OldPwd  string `json:"old_password"`
	NewPwd  string `json:"new_password"`
	ConfPwd string `json:"confirm_password"`
}
