package repositories

import (
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
func GetPrice(token *entities.Users, r *entities.GetPriceReq) (*entities.GetPriceRes, string, string, bool) {
	if r.DestQty == 0 {
		return nil, "99", "destination_qty is required", false
	}

	if len(r.Visit) <= 0 {
		return nil, "99", "tarif is required", false
	}

	var agent entities.AgentModel

	if err := db.DB[1].Select(`agent_group_id`).Where("agent_id = ?", token.Typeid).Find(&agent).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Agent not found (" + err.Error() + ")", false
	}

	if agent.Agent_group_id == 0 {
		return nil, "03", "Agent group not configured yet, please contact the Administrator", false
	}

	var curr entities.CurrencyModel

	if err := db.DB[1].Last(&curr).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "05", "Currency not found (" + err.Error() + ")", false
	}

	var vis []entities.VisitRes
	var priceTotal, discTotal float32

	for _, visit := range r.Visit {
		var trfRes []entities.TrfPriceRes

		for _, trf := range visit.Trf {
			var tarif entities.TariffModel

			if err := db.DB[1].Where("trf_id = ?", trf.ID).Find(&tarif).Error; gorm.IsRecordNotFoundError(err) {
				return nil, "04", "Tarif not found (" + err.Error() + ")", false
			}

			var dD entities.DiscountMultiModel
			var dA entities.DiscountMultiModel
			dDStatus := true
			dAStatus := true

			if err := db.DB[1].Select(`discm_value`).Where(`discm_type = 'MULTIDESTINATION'
									and ? >= discm_start_date
									and ? <= discm_end_date 
									and discm_destination = ?
									and deleted_at is NULL`, visit.VisitDate, visit.VisitDate, r.DestQty).Last(&dD).Error; gorm.IsRecordNotFoundError(err) {
				dDStatus = false
			}

			if err := db.DB[1].Select(`discm_value`).Where(`discm_type = 'AGENT'
									and ? >= discm_start_date
									and ? <= discm_end_date 
									and discm_destination = ?
									and discm_group_agent_id = ?
									and deleted_at is NULL`, visit.VisitDate, visit.VisitDate, r.DestQty, agent.Agent_group_id).Last(&dA).Error; gorm.IsRecordNotFoundError(err) {
				dAStatus = false
			}

			trfVal := tarif.Trf_value

			if tarif.Trf_currency_code == "USD" {
				trfVal = curr.Curr_rate * tarif.Trf_value
			}
			discDes := (dD.Discm_value / 100) * trfVal
			discAg := (dA.Discm_value / 100) * (trfVal - discDes)
			totDes := discDes * float32(trf.Qty)
			totAg := discAg * float32(trf.Qty)
			totVal := trfVal * float32(trf.Qty)

			priceTotal += totVal
			discTotal += (totDes + totAg)

			tmpPrice := entities.TrfPriceRes{
				TrfID:              trf.ID,
				TrfName:            tarif.Trf_name,
				TrfCode:            tarif.Trf_code,
				TrfValue:           trfVal,
				Currency:           tarif.Trf_currency_code,
				Qty:                trf.Qty,
				DiscDStatus:        dDStatus,
				DiscDPercent:       dD.Discm_value,
				DiscDestination:    discDes,
				DiscAStatus:        dAStatus,
				DiscAPercent:       dA.Discm_value,
				DiscAgent:          discAg,
				TotValue:           totVal,
				TotDiscDestination: totDes,
				TotDiscAgent:       totAg,
				TotPrice:           totVal - totDes - totAg,
			}

			trfRes = append(trfRes, tmpPrice)
		}

		tmpVisit := entities.VisitRes{
			VisitDate: visit.VisitDate,
			Trf:       trfRes,
		}

		vis = append(vis, tmpVisit)
	}

	price := entities.GetPriceRes{
		Visit:        vis,
		PriceTotal:   priceTotal,
		Discount:     discTotal,
		PaymentTotal: priceTotal - discTotal,
	}

	return &price, "01", "Success get price", true
}
