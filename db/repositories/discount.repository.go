package repositories

import (
	"strconv"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"

	"github.com/jinzhu/gorm"
)

//GetDiscountMulti : select data agent
func GetDiscountMulti(token *entities.Users, discountType string) (*[]entities.DiscountMultiModel, string, string, bool) {
	var agent entities.AgentModel

	if err := db.DB[1].Select(`agent_group_id`).Where("agent_id = ?", token.Typeid).Find(&agent).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Agent not found (" + err.Error() + ")", false
	}

	if agent.Agent_group_id == 0 {
		return nil, "03", "Agent group not configured yet, please contact the Administrator", false
	}

	var discount []entities.DiscountMultiModel

	if discountType == "AGENT" {
		if err := db.DB[1].Select(`discm_name, discm_start_date::text, discm_end_date::text, discm_destination, discm_value`).Where(`discm_group_agent_id = ?
								and discm_type = ?
								and current_date >= discm_start_date
								and current_date <= discm_end_date 
								and deleted_at is NULL`, agent.Agent_group_id, discountType).Order("discm_name, discm_destination").Find(&discount).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "04", "Discount not found (" + err.Error() + ")", false
		}
	} else {
		if err := db.DB[1].Select(`discm_name, discm_start_date::text, discm_end_date::text, discm_destination, discm_value`).Where(`discm_type = ?
								and current_date >= discm_start_date
								and current_date <= discm_end_date 
								and deleted_at is NULL`, discountType).Order("discm_name, discm_destination").Find(&discount).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "04", "Discount not found (" + err.Error() + ")", false
		}
	}

	return &discount, "01", "Get data discount success", true
}

//GetPrice : get price data
func GetPrice(token *entities.Users, r *entities.GetPriceReq) (*[]entities.GetPriceRes, string, string, bool) {
	var price []entities.GetPriceRes

	for i, trf := range r.Trf {
		tmpPrice := entities.GetPriceRes{
			TrfID:              trf.ID,
			TrfName:            "Coba " + strconv.Itoa(i),
			TrfCode:            "abcd-ex",
			TrfValue:           666,
			Qty:                trf.Qty,
			DiscDStatus:        true,
			DiscDPercent:       5,
			DiscDestination:    2500,
			DiscAStatus:        true,
			DiscAPercent:       5,
			DiscAgent:          2600,
			TotValue:           100000,
			TotDiscDestination: 2500,
			TotDiscAgent:       2600,
			TotPrice:           94100,
		}

		price = append(price, tmpPrice)
	}

	return &price, "01", "Success get price", true
}
