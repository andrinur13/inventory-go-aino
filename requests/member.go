package requests

type MemberRequest struct {
	DataMember []Member `json:"data_member"`
}

type Member struct {
	SerialNumber     string `json:"serial_number"`
	ExpiredDate      string `json:"expired_date"`
	Name             string `json:"member_name"`
	RegistrationDate string `json:"registration_date"`
	TypeCode         int    `json:"type_code"`
	Type             string `json:"type"`
	DayLimit         int    `json:"day_limit"`
}
