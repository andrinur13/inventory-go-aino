package repositories

import (
	"encoding/json"
	"time"
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
								agent_extras ->> 'email' as email,
								agent_extras ->> 'npwp' as npwp,
								agent_extras ->> 'pic_name' as pic_name`).Where("deleted_at is null").Order("agent_name").Find(&agents).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "No agent data found (" + err.Error() + ")", false
	}

	var dataAgent []entities.AgentModel

	for _, agent := range agents {
		extras := entities.AgentExtras{
			NoID:       agent.No_id,
			AddrDetail: agent.Agent_address_detail,
			PicName:    agent.Pic_name,
			Telp:       agent.Telp,
			Email:      agent.Email,
			Npwp:      	agent.Npwp,
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

//GetAgent : select data agent
func GetDetailAgent(token *entities.Users) (*[]entities.AgentModel, string, string, bool) {
	var agents []entities.AgentModel

	if err := db.DB[0].Select(`agent_id, agent_name, agent_address,
								agent_extras ->> 'agent_address_detail' as agent_address_detail,
								agent_extras ->> 'telp' as telp,
								agent_extras ->> 'no_id' as no_id,
								agent_extras ->> 'email' as email,
								agent_extras ->> 'npwp' as npwp,
								agent_extras ->> 'pic_name' as pic_name`).Where("agent_id = ? AND deleted_at is null", token.Typeid).Order("agent_name").Find(&agents).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "No agent data found (" + err.Error() + ")", false
	}

	var dataAgent []entities.AgentModel

	for _, agent := range agents {
		extras := entities.AgentExtras{
			NoID:       agent.No_id,
			AddrDetail: agent.Agent_address_detail,
			PicName:    agent.Pic_name,
			Telp:       agent.Telp,
			Email:      agent.Email,
			Npwp:				agent.Npwp,
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

func UpdateProfileAgent(token *entities.Users, r *entities.AgentReq) (map[string]interface{}, string, string, bool) {
	if r.Agent == "" {
		return nil, "99", "Agent cant't be empty", false
	}

	if r.Address == "" {
		return nil, "99", "Address cant't be empty", false
	}

	if r.Extras.AddrDetail == "" {
		return nil, "99", "Addr detail cant't be empty", false
	}

	if r.Extras.NoID == "" {
		return nil, "99", "No ID cant't be empty", false
	}

	if r.Extras.PicName == "" {
		return nil, "99", "Pic name cant't be empty", false
	}

	if r.Extras.Telp == "" {
		return nil, "99", "Telp number cant't be empty", false
	}

	if r.Extras.Email == "" {
		return nil, "99", "E-mail cant't be empty", false
	}

	if r.Extras.Npwp == "" {
		return nil, "99", "Npwp cant't be empty", false
	}

	rExt, err := json.Marshal(&r.Extras)
	if err != nil {
		return nil, "99", "Failed to parse json key contact (" + err.Error() + ")", false
	}

	jExt := string(rExt)

	var agent entities.AgentModel

	if r.Group != 0 {
		agent = entities.AgentModel{
			Agent_name:     r.Agent,
			Agent_address:  r.Address,
			Agent_group_id: r.Group,
			Agent_extras:   jExt,
		}
	} else {
		agent = entities.AgentModel{
			Agent_name:    r.Agent,
			Agent_address: r.Address,
			Agent_extras:  jExt,
		}
	}

	var checkAgent []entities.AgentModel

	db.DB[1].Where("deleted_at is null and agent_id = ?", token.Typeid).Find(&checkAgent).Update(agent)

	return nil, "01", "Agent successfully updated", true
}

//InsertAgent : insert data user
func InsertAgent(token *entities.Users, r *entities.AgentReq) (map[string]interface{}, string, string, bool) {
	if r.Agent == "" {
		return nil, "99", "Agent cant't be empty", false
	}

	if r.Address == "" {
		return nil, "99", "Address cant't be empty", false
	}

	if r.Extras.AddrDetail == "" {
		return nil, "99", "Addr detail cant't be empty", false
	}

	if r.Extras.NoID == "" {
		return nil, "99", "No ID cant't be empty", false
	}

	if r.Extras.PicName == "" {
		return nil, "99", "Pic name cant't be empty", false
	}

	if r.Extras.Telp == "" {
		return nil, "99", "Telp number cant't be empty", false
	}

	if r.Extras.Email == "" {
		return nil, "99", "E-mail cant't be empty", false
	}

	if r.Extras.Npwp == "" {
		return nil, "99", "NPWP cant't be empty", false
	}

	var checkAgent []entities.AgentModel

	db.DB[1].Where("deleted_at is null and agent_name = ?", r.Agent).Find(&checkAgent)

	if len(checkAgent) > 0 {
		return nil, "02", "Agent with specified data already registered", false
	}

	rExt, err := json.Marshal(&r.Extras)
	if err != nil {
		return nil, "99", "Failed to parse json key contact (" + err.Error() + ")", false
	}

	jExt := string(rExt)

	var agent entities.AgentModel

	if r.Group != 0 {
		agent = entities.AgentModel{
			Agent_name:     r.Agent,
			Agent_address:  r.Address,
			Agent_group_id: r.Group,
			Agent_extras:   jExt,
			Created_at:     time.Now().Format("2006-01-02 15:04:05"),
		}
	} else {
		agent = entities.AgentModel{
			Agent_name:    r.Agent,
			Agent_address: r.Address,
			Agent_extras:  jExt,
			Created_at:    time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	db.DB[1].NewRecord(agent)

	if err := db.DB[1].Create(&agent).Error; err != nil {
		return nil, "03", "Error when inserting agent data (" + err.Error() + ")", false
	}

	return map[string]interface{}{
		"agent":   r.Agent,
		"address": r.Address,
		"group":   r.Group,
		"contact": r.Extras,
	}, "01", "Agent registration success", true
}
