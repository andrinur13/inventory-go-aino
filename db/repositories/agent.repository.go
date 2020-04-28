package repositories

import (
	"twc-ota-api/db"
	"twc-ota-api/db/entities"

	"github.com/jinzhu/gorm"
)

//GetAgent : select data agent
func GetAgent() (*[]entities.AgentModel, string, string, bool) {
	var agents []entities.AgentModel

	if err := db.DB[0].Select(`agent_id, agent_name, agent_address,
								agent_extras ->> 'agent_address_detail' as agent_address_detail,
								agent_extras ->> 'telp' as telp,
								agent_extras ->> 'no_id' as no_id,
								agent_extras ->> 'pic_name' as pic_name`).Order("agent_name").Find(&agents).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "No agent data found (" + err.Error() + ")", false
	}

	var dataAgent []entities.AgentModel

	for _, agent := range agents {
		extras := entities.AgentExtras{
			NoID:       agent.No_id,
			AddrDetail: agent.Agent_address_detail,
			PicName:    agent.Pic_name,
			Telp:       agent.Telp,
		}

		tmpAgent := entities.AgentModel{
			Agent_id:      agent.Agent_id,
			Agent_address: agent.Agent_address,
			Agent_name:    agent.Agent_name,
			AgentExtras:   &extras,
		}

		dataAgent = append(dataAgent, tmpAgent)
	}

	return &dataAgent, "01", "Get data agent success", true
}
