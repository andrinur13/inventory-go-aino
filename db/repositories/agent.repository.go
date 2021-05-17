package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"
	"math"
	"reflect"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"

	"github.com/jinzhu/gorm"
)

//GetAgent : select data agent
func GetAgent() (*[]entities.AgentModel, string, string, bool) {
	var agents []entities.AgentModel

	if err := db.DB[1].Select(`agent_id, agent_name, agent_address, group_agent_name,
								agent_extras ->> 'agent_address_detail' as agent_address_detail,
								agent_extras ->> 'telp' as telp,
								agent_extras ->> 'no_id' as no_id,
								agent_extras ->> 'email' as email,
								agent_extras ->> 'npwp' as npwp,
								agent_extras ->> 'pic_name' as pic_name`).Where("master_agents.deleted_at is null").Joins("inner join master_agents_group on group_agent_id = master_agents.agent_group_id").Order("agent_name").Find(&agents).Error; gorm.IsRecordNotFoundError(err) {
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
			Agent_group_id: 	agent.Agent_group_id,
			AgentExtras:   &extras,
		}

		dataAgent = append(dataAgent, tmpAgent)
	}

	return &dataAgent, "01", "Get data agent success", true
}

//GetAgent : select data agent
func GetDetailAgent(token *entities.Users) (*[]entities.AgentModel, string, string, bool) {
	var agents []entities.AgentModel

	err := db.DB[1].Select(`agent_id, agent_name, agent_address, agent_group_id, group_agent_name,
								case
									when ((agent_extras -> 'image_url') isnull) then 'b2bm/agent/default_agent.png'
									when ((agent_extras ->> 'image_url') = '') then 'b2bm/agent/default_agent.png'
									else agent_extras ->> 'image_url'
								end as agent_image_url,
								agent_extras ->> 'agent_address_detail' as agent_address_detail,
								agent_extras ->> 'telp' as telp,
								agent_extras ->> 'no_id' as no_id,
								agent_extras ->> 'email' as email,
								agent_extras ->> 'npwp' as npwp,
								agent_extras ->> 'pic_name' as pic_name`).Where("master_agents.agent_id = ? AND master_agents.deleted_at is null", token.Typeid).Joins("inner join master_agents_group on group_agent_id = master_agents.agent_group_id").Order("agent_name").Find(&agents).Error;

	//If Connection Refused
	if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
		fmt.Printf("%v \n", err.Error())
			for i := 0; i<4; i++ {
				err = db.DB[1].Select(`agent_id, agent_name, agent_address, agent_group_id, group_agent_name,
								case
									when ((agent_extras -> 'image_url') isnull) then 'b2bm/agent/default_agent.png'
									when ((agent_extras ->> 'image_url') = '') then 'b2bm/agent/default_agent.png'
									else agent_extras ->> 'image_url'
								end as agent_image_url,
								agent_extras ->> 'agent_address_detail' as agent_address_detail,
								agent_extras ->> 'telp' as telp,
								agent_extras ->> 'no_id' as no_id,
								agent_extras ->> 'email' as email,
								agent_extras ->> 'npwp' as npwp,
								agent_extras ->> 'pic_name' as pic_name`).Where("master_agents.agent_id = ? AND master_agents.deleted_at is null", token.Typeid).Joins("inner join master_agents_group on group_agent_id = master_agents.agent_group_id").Order("agent_name").Find(&agents).Error;
				if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
					fmt.Printf("Hitback(%d)%v \n", i, err)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	if gorm.IsRecordNotFoundError(err) {
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
			Npwp:		agent.Npwp,
		}

		tmpAgent := entities.AgentModel{
			Agent_id:       	agent.Agent_id,
			Group_agent_name: 	agent.Group_agent_name,
			Agent_address:  	agent.Agent_address,
			Agent_name:     	agent.Agent_name,
			Agent_group_id: 	agent.Agent_group_id,
			AgentExtras:    	&extras,
			AgentEmail:			token.Email,
			AgentUsername:		token.Name,
			Agent_image_url:	agent.Agent_image_url,
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

	//If Connection refused not yet
	if erro := db.DB[1].Where("deleted_at is null and agent_id = ?", token.Typeid).Find(&checkAgent).Error; (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError") {
		fmt.Printf("%v \n", erro.Error())
		fmt.Printf("%v \n", reflect.TypeOf(erro).String())
			for i := 0; i<4; i++ {
				erro = db.DB[1].Where("deleted_at is null and agent_id = ?", token.Typeid).Find(&checkAgent).Error;
				if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError") {
					fmt.Printf("Hitback(%d)%v \n", i, erro)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}
	
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

//GetInbox : inbox notification
func GetInboxNotification(token *entities.Users, typeNotif string, page int, size int) (*[]entities.InboxNotificationModel, string, string, bool, int, int, int) {
	var inbox []entities.InboxNotificationModel
	var countInbox []entities.InboxNotificationModel

	offset := (page - 1) * size
	limit := size

	now := time.Now()
	
	var typeN int;
	if (typeNotif != ""){
		typeNo, err := strconv.Atoi(typeNotif)
		if err != nil {
			// handle error
			fmt.Println(err)
			return nil, "99", "Failed to parse type", false, 0, 0, 0
		}
		typeN = typeNo
	}else{
		typeN = 0
	}

	if (typeN == 1) || (typeN == 2){
		err := db.DB[0].Select(`inbox_id`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_type = ? and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, typeN, now, now).Find(&countInbox).Error;

		//If Connection refused
		if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
			fmt.Printf("%v \n", err.Error())
				for i := 0; i<4; i++ {
					err = db.DB[0].Select(`inbox_id`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_type = ? and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, typeN, now, now).Find(&countInbox).Error;
					if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError") {
						fmt.Printf("Hitback(%d)%v \n", i, err)
						time.Sleep(3 * time.Second)
						continue
					}
					break
				}
			if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
				return nil, "502", "Connection has a problem", false, 0, 0, 0
			}
		}

		if len(countInbox) == 0 {
			return nil, "60", "Inbox not found", false, 0, 0, 0
		}

		if gorm.IsRecordNotFoundError(err) {
			return nil, "60", "Inbox not found (" + err.Error() + ")", false, 0, 0, 0
		}

		if err := db.DB[0].Select(`inbox_id,
									agent_name as agent_name,
									group_agent_name as agent_group_name,
									inbox_created_at as inbox_created_at,
									case
										when inbox_image_url isnull then 'b2bm/inbox/default_inbox.jpg'
										else inbox_image_url
									end as inbox_image_url,
									inbox_title as inbox_title,
									inbox_subtitle as inbox_short_desc,
									inbox_desc as inbox_full_desc
									`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_type = ? and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, typeN, now, now).Joins("inner join master_agents on agent_id = inbox_agent_id").Joins("inner join master_agents_group on group_agent_id = inbox_group_agent_id").Limit(limit).Offset(offset).Find(&inbox).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "60", "Inbox not found (" + err.Error() + ")", false, 0, 0, 0
		}

	}else {
		err := db.DB[0].Select(`inbox_id`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, now, now).Find(&countInbox).Error;

		//If Connection refused
		if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
			fmt.Printf("%v \n", err.Error())
				for i := 0; i<4; i++ {
					err = db.DB[0].Select(`inbox_id`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, now, now).Find(&countInbox).Error;
					if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError") {
						fmt.Printf("Hitback(%d)%v \n", i, err)
						time.Sleep(3 * time.Second)
						continue
					}
					break
				}
			if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
				return nil, "502", "Connection has a problem", false, 0, 0, 0
			}
		}

		if len(countInbox) == 0 {
			return nil, "60", "Inbox not found", false, 0, 0, 0
		}

		if gorm.IsRecordNotFoundError(err) {
			return nil, "60", "Inbox not found (" + err.Error() + ")", false, 0, 0, 0
		}

		if err := db.DB[0].Select(`inbox_id,
									agent_name as agent_name,
									group_agent_name as agent_group_name,
									inbox_created_at as inbox_created_at,
									case
										when inbox_image_url isnull then 'b2bm/inbox/default_inbox.jpg'
										else inbox_image_url
									end as inbox_image_url,
									inbox_title as inbox_title,
									inbox_subtitle as inbox_short_desc,
									inbox_desc as inbox_full_desc
									`).Where("(inbox_agent_id = ? or inbox_agent_id = 0) and inbox_show_start_date <= ? and (inbox_show_end_date >= ? or inbox_show_end_date isnull)", token.Typeid, now, now).Joins("inner join master_agents on agent_id = inbox_agent_id").Joins("inner join master_agents_group on group_agent_id = inbox_group_agent_id").Limit(limit).Offset(offset).Find(&inbox).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "60", "Inbox not found (" + err.Error() + ")", false, 0, 0, 0
		}

	}

	var dataInbox []entities.InboxNotificationModel

	for _, data := range inbox {

		createdat, _ := time.Parse(time.RFC3339, data.Inbox_created_at)

		tmpInbox := entities.InboxNotificationModel{
			Inbox_id:      		data.Inbox_id,
			Agent_name: 		data.Agent_name,
			Agent_group_name:  	data.Agent_group_name,
			Inbox_created_at:   createdat.Format("02 Jan 2006"),
			Inbox_image_url:  	data.Inbox_image_url,
			Inbox_title:  		data.Inbox_title,
			Inbox_short_desc: 	data.Inbox_short_desc,
			Inbox_full_desc:  	data.Inbox_full_desc,
		}

		dataInbox = append(dataInbox, tmpInbox)
	}

	totalData := len(countInbox)
	currentData := len(inbox)
	rawPages := float64(totalData) / float64(size)
	totalPages := int(math.Ceil(rawPages))

	return &dataInbox, "01", "Success get inbox list.", true, totalData, totalPages, currentData
}
