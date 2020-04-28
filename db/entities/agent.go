package entities

type AgentModel struct {
	Agent_id             int          `json:"agent_id"`
	Agent_name           string       `json:"agent_name"`
	Agent_address        string       `json:"agent_address"`
	Agent_address_detail string       `json:"-"`
	Telp                 string       `json:"-"`
	No_id                string       `json:"-"`
	Pic_name             string       `json:"-"`
	AgentExtras          *AgentExtras `json:"agent_extras"`
}

func (AgentModel) TableName() string {
	return "master_agents"
}

type AgentExtras struct {
	NoID       string `json:"no_id"`
	AddrDetail string `json:"address_detail"`
	Telp       string `json:"telp"`
	PicName    string `json:"pic_name"`
}
