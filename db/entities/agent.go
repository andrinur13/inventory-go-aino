package entities

type AgentModel struct {
	Agent_id             int          `json:"agent_id" gorm:"primary_key"`
	Agent_name           string       `json:"agent_name"`
	Agent_address        string       `json:"agent_address"`
	Agent_address_detail string       `json:"-" gorm:"-"`
	Telp                 string       `json:"-" gorm:"-"`
	No_id                string       `json:"-" gorm:"-"`
	Pic_name             string       `json:"-" gorm:"-"`
	Email                string       `json:"-" gorm:"-"`
	Npwp                 string       `json:"-" gorm:"-"`
	Agent_group_id       int          `json:"-" gorm:"DEFAULT null"`
	Agent_extras         string       `json:"-"`
	AgentExtras          *AgentExtras `json:"agent_extras" gorm:"-"`
	Created_at           string       `json:"-"`
}

func (AgentModel) TableName() string {
	return "master_agents"
}

type AgentExtras struct {
	NoID       string `json:"no_id"`
	AddrDetail string `json:"agent_address_detail"`
	Telp       string `json:"telp"`
	PicName    string `json:"pic_name"`
	Email      string `json:"email"`
	Npwp       string `json:"npwp"`
}

type AgentReq struct {
	Agent   string      `json:"agent"`
	Address string      `json:"address"`
	Group   int         `json:"group"`
	Extras  AgentExtras `json:"contact"`
}
