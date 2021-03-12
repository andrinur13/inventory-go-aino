package repositories

import (
	"fmt"
	"math"
	"strings"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/utils/helper"

	"github.com/jinzhu/gorm"
)

var qrPrefix = "AINO"
var prefixLen = len(qrPrefix)
var weekday = strings.ToUpper("%" + time.Now().Format("Mon") + "%")

// GetTicket : get ticket list for device
func GetTicket(r interface{}, token *entities.Users) (map[string]interface{}, string, string, bool) {
	var masterTicket []entities.MasterTicket
	var mTickets []entities.MasterTicket

	// query := `SELECT
	// 			COALESCE(cast(a.trf_id as varchar(255)), '') as trf_id,
	// 			COALESCE(a.trf_code, '') as trf_code,
	// 			COALESCE(cast(a.trf_group_id as varchar(255)), '') as trf_group_id,
	// 			COALESCE(d.group_name , '') as trf_group_name,
	// 			COALESCE(f.trfftype_name, '') as trf_trftype,
	// 			COALESCE(a.trf_name, '') as trf_name,
	// 			COALESCE(a.trf_value::text, '') as trf_value,
	// 			COALESCE(a.trf_start_date::text, '') as trf_start_date,
	// 			COALESCE(a.trf_end_date::text, '') as trf_end_date,
	// 			COALESCE(a.trf_priority::text, '') as trf_priority,
	// 			COALESCE(a.trf_release::text, '') as trf_release,
	// 			COALESCE(a.trf_currency_code, '') as trf_currency_code,
	// 			COALESCE(a.trf_qty, '1') as trf_qty,
	// 			a.trf_agent_id,
	// 			COALESCE(cast(trf_condition->>'day' as text), '') as day,
	// 			COALESCE(cast(trf_condition->>'beginTime' as text), '') as begin_time,
	// 			COALESCE(cast(trf_condition->>'endTime' as text), '') as end_time,
	// 			COALESCE(cast(trf_condition->>'cardType' as text), '') as card_type,
	// 			COALESCE(cast(trf_condition->>'expiredQr' as text), '') as expired_qr
	// 			from public.master_tariff a
	// 			join public.master_group d on a.trf_group_id = d.group_id
	// 			left join public.master_tariff_type f on a.trf_trfftype_id = f.trfftype_id
	// 			left join public.master_tariff_has_machine i on a.trf_id = i.trf_id
	// 			where a.deleted_at is null
	// 			and a.trf_agent_id = ?`
	query := `SELECT 
				COALESCE(cast(a.trf_id as varchar(255)), '') as trf_id,
				COALESCE(a.trf_name, '') as trf_name,
				COALESCE(a.trf_value::text, '') as trf_value,
				COALESCE(a.trf_currency_code, '') as trf_currency_code
				from public.master_tariff a
				join public.master_group d on a.trf_group_id = d.group_id
				where a.deleted_at is null
				and a.trf_agent_id = ?`

	mbmid := r.(map[string]interface{})["mbmid"]

	if mbmid == nil || mbmid == "" {
		query += ` order by a.trf_id DESC;`
		db.DB[0].Raw(query, token.Typeid).Scan(&masterTicket)
	} else {
		query += ` and (d.group_mid = ?) order by a.trf_id DESC;`
		db.DB[0].Raw(query, token.Typeid, mbmid).Scan(&masterTicket)
	}

	if len(masterTicket) == 0 {
		return nil, "02", "Fare not found", false
	}

	var curr entities.CurrencyModel

	if err := db.DB[1].Last(&curr).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "05", "Currency not found (" + err.Error() + ")", false
	}

	for _, item := range masterTicket {
		trfVal := item.Trf_value
		trfLabel := "Rp. " + helper.RenderFloat("#.###,", float64(trfVal))

		if item.Trf_currency_code == "USD" {
			trfVal = curr.Curr_rate * item.Trf_value
			trfLabel = "USD $" + fmt.Sprintf("%g", item.Trf_value) + " | " + "Rp. " + helper.RenderFloat("#.###,", float64(trfVal))
		}

		mTicket := entities.MasterTicket{
			Trf_id:            item.Trf_id,
			Trf_name:          item.Trf_name,
			Trf_value:         trfVal,
			Trf_label:         trfLabel,
			Trf_currency_code: item.Trf_currency_code,
		}

		mTickets = append(mTickets, mTicket)
	}

	var site entities.GrupModel

	if err := db.DB[1].Select("*, COALESCE(cast(group_extras ->>'estimate' as text), '') as group_estimate").Where("group_mid = ?", mbmid).Find(&site).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "04", "Site not found (" + err.Error() + ")", false
	}

	return map[string]interface{}{
		"site_id":       site.Group_id,
		"mmid":          site.Group_mid,
		"site_name":     site.Group_name,
		"site_images":   site.Group_logo,
		"site_duration": site.Group_estimate,
		"ticket_list":   mTickets,
	}, "01", "Success get ticket list", true
}

//SelectCluster : select data cluster
func SelectCluster(token *entities.Users, nationality string) (*[]entities.Cluster, string, string, bool) {

	var cluster []entities.GrupModel

	if err := db.DB[1].Select("*, coalesce(cast(group_extras ->> 'detail' as text), '') as description").Where("depth = 2 AND cast(group_extras ->> 'type' as text) = 'TP' AND deleted_at is NULL").Find(&cluster).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Cluster not found (" + err.Error() + ")", false
	}

	var resp []entities.Cluster

	for _, data := range cluster {
		var sites []entities.GrupModel

		if err := db.DB[1].Select(`*, COALESCE(cast(group_extras ->>'estimate' as text), '') as group_estimate, coalesce(cast(group_extras ->> 'latitude' as text), '') as lat,
									coalesce(cast(group_extras ->> 'longitude' as text), '') as long`).Where("depth = 3 AND cast(group_extras ->> 'type' as text) = 'TP' AND parent_id = ? AND deleted_at is NULL", data.Group_id).Find(&sites).Error; err != nil {
			return nil, "03", "Error when fetching site data (" + err.Error() + ")", false
		}

		var siteResp []entities.Site

		for _, site := range sites {
			var adultTrf []entities.SiteTrfModel

			if nationality == "" {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAD%' OR trf_code ilike '%TVDW%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			} else if nationality == "96" {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVDW%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			} else {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAD%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			}

			var dataAdult []entities.SiteTrfModel

			for _, aTrf := range adultTrf {
				var getTicket []entities.SiteTrfModel

				if err := db.DB[1].Select("mtick_name").Where("trf_id = ?", aTrf.Trf_id).Joins(`inner join master_tariffdet on trf_id = trfdet_trf_id 
								inner join master_ticket on trfdet_mtick_id = mtick_id`).Find(&getTicket).Error; err != nil {
					return nil, "05", err.Error(), false
				}

				var ticks string

				for _, tick := range getTicket {
					ticks += tick.Mtick_name + ", "
				}

				r := []rune(ticks)

				if len(r) > 2 {
					ticks = string(r[:len(r)-2])
				}

				tmpAdult := entities.SiteTrfModel{
					Trf_id:            aTrf.Trf_id,
					Trf_code:          aTrf.Trf_code,
					Trf_name:          aTrf.Trf_name,
					Trf_currency_code: aTrf.Trf_currency_code,
					Trfftype_name:     aTrf.Trfftype_name,
					Trf_value:         aTrf.Trf_value,
					Mtick_name:        ticks,
				}

				dataAdult = append(dataAdult, tmpAdult)
			}

			var childTrf []entities.SiteTrfModel

			if nationality == "" {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAN%' OR trf_code ilike '%TVCH%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			} else if nationality == "96" {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAN%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			} else {
				if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVCH%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
					return nil, "04", err.Error(), false
				}
			}

			var dataChild []entities.SiteTrfModel

			for _, cTrf := range childTrf {
				var getTicket []entities.SiteTrfModel

				if err := db.DB[1].Select("mtick_name").Where("trf_id = ?", cTrf.Trf_id).Joins(`inner join master_tariffdet on trf_id = trfdet_trf_id 
								inner join master_ticket on trfdet_mtick_id = mtick_id`).Find(&getTicket).Error; err != nil {
					return nil, "05", err.Error(), false
				}

				var ticks string

				for _, tick := range getTicket {
					ticks += tick.Mtick_name + ", "
				}

				r := []rune(ticks)

				if len(r) > 2 {
					ticks = string(r[:len(r)-2])
				}

				tmpChild := entities.SiteTrfModel{
					Trf_id:            cTrf.Trf_id,
					Trf_code:          cTrf.Trf_code,
					Trf_name:          cTrf.Trf_name,
					Trf_currency_code: cTrf.Trf_currency_code,
					Trfftype_name:     cTrf.Trfftype_name,
					Trf_value:         cTrf.Trf_value,
					Mtick_name:        ticks,
				}

				dataChild = append(dataChild, tmpChild)
			}

			tmpSite := entities.Site{
				SiteID:        site.Group_id,
				SiteMID:       site.Group_mid,
				SiteName:      site.Group_name,
				SiteLogo:      site.Group_logo,
				SiteEstimated: site.Group_estimate,
				SiteLat:       site.Lat,
				SiteLong:      site.Long,
				Trf: entities.SiteTrf{
					Adult: dataAdult,
					Child: dataChild,
				},
			}

			siteResp = append(siteResp, tmpSite)
		}

		tmpResp := entities.Cluster{
			ClusterID:          data.Group_id,
			ClusterMID:         data.Group_mid,
			ClusterName:        data.Group_name,
			ClusterLogo:        data.Group_logo,
			Site:               siteResp,
			ClusterDescription: data.Description,
		}

		resp = append(resp, tmpResp)
	}

	return &resp, "01", "Success get cluster data", true
}

//GetSite : select data detail site
func GetSite(token *entities.Users, nationality string, siteID string) (*entities.SiteDetail, string, string, bool) {

	if siteID == "" {
		return nil, "02", "site_id param is required", false
	}

	var site entities.GroupSiteModel

	if err := db.DB[1].Select(`group_id, group_name, group_mid, group_logo,
								coalesce(cast(group_extras ->> 'open' as text), '') as "open",
								coalesce(cast(group_extras ->> 'close' as text), '') as "close",
								coalesce(cast(group_extras ->> 'estimate' as text), '') as estimate,
								coalesce(cast(group_extras ->> 'detail' as text), '') as detail,
								coalesce(cast(group_extras ->> 'address' as text), '') as address,
								coalesce(cast(group_extras ->> 'latitude' as text), '') as lat,
								coalesce(cast(group_extras ->> 'longitude' as text), '') as long`).Where("group_id = ? AND depth = 3 AND cast(group_extras ->> 'type' as text) = 'TP' AND deleted_at is NULL", siteID).First(&site).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "03", "Site not found (" + err.Error() + ")", false
	}

	var adultTrf []entities.SiteTrfModel

	if nationality == "" {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAD%' OR trf_code ilike '%TVDW%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	} else if nationality == "96" {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVDW%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	} else {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAD%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&adultTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	}

	var dataAdult []entities.SiteTrfModel

	for _, aTrf := range adultTrf {
		var getTicket []entities.SiteTrfModel

		if err := db.DB[1].Select("mtick_name").Where("trf_id = ?", aTrf.Trf_id).Joins(`inner join master_tariffdet on trf_id = trfdet_trf_id 
								inner join master_ticket on trfdet_mtick_id = mtick_id`).Find(&getTicket).Error; err != nil {
			return nil, "05", err.Error(), false
		}

		var ticks string

		for _, tick := range getTicket {
			ticks += tick.Mtick_name + ", "
		}

		r := []rune(ticks)

		if len(r) > 2 {
			ticks = string(r[:len(r)-2])
		}

		tmpAdult := entities.SiteTrfModel{
			Trf_id:            aTrf.Trf_id,
			Trf_code:          aTrf.Trf_code,
			Trf_name:          aTrf.Trf_name,
			Trf_currency_code: aTrf.Trf_currency_code,
			Trfftype_name:     aTrf.Trfftype_name,
			Trf_value:         aTrf.Trf_value,
			Mtick_name:        ticks,
		}

		dataAdult = append(dataAdult, tmpAdult)
	}

	var childTrf []entities.SiteTrfModel

	if nationality == "" {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAN%' OR trf_code ilike '%TVCH%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	} else if nationality == "96" {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVAN%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	} else {
		if err := db.DB[1].Select(`trf_id,
											trf_name,
											trf_code,
											CASE
												WHEN trf_currency_code != 'IDR' THEN (trf_value::NUMERIC*curr.curr_rate::NUMERIC)
												ELSE trf_value
											END AS trf_value,
											trf_currency_code, trfftype_name`).Where(`trf_group_id = ?
											AND (trf_code ilike '%TVCH%')
											AND (trf_end_date IS NULL OR trf_end_date >= LOCALTIMESTAMP)
											AND trf_agent_id = ?
											AND deleted_at is NULL
								`, site.Group_id, token.Typeid).Joins(`inner join master_tariff_type ON trfftype_id = trf_trfftype_id LEFT JOIN
															(SELECT *
															FROM currency c
															ORDER BY c.created_at DESC
															LIMIT 1) curr ON curr.created_at <= master_tariff.created_at
														OR curr.created_at >= master_tariff.created_at`).Find(&childTrf).Error; err != nil {
			return nil, "04", err.Error(), false
		}
	}

	var dataChild []entities.SiteTrfModel

	for _, cTrf := range childTrf {
		var getTicket []entities.SiteTrfModel

		if err := db.DB[1].Select("mtick_name").Where("trf_id = ?", cTrf.Trf_id).Joins(`inner join master_tariffdet on trf_id = trfdet_trf_id 
								inner join master_ticket on trfdet_mtick_id = mtick_id`).Find(&getTicket).Error; err != nil {
			return nil, "05", err.Error(), false
		}

		var ticks string

		for _, tick := range getTicket {
			ticks += tick.Mtick_name + ", "
		}

		r := []rune(ticks)

		if len(r) > 2 {
			ticks = string(r[:len(r)-2])
		}

		tmpChild := entities.SiteTrfModel{
			Trf_id:            cTrf.Trf_id,
			Trf_code:          cTrf.Trf_code,
			Trf_name:          cTrf.Trf_name,
			Trf_currency_code: cTrf.Trf_currency_code,
			Trfftype_name:     cTrf.Trfftype_name,
			Trf_value:         cTrf.Trf_value,
			Mtick_name:        ticks,
		}

		dataChild = append(dataChild, tmpChild)
	}

	resp := entities.SiteDetail{
		Group_id:   site.Group_id,
		Group_mid:  site.Group_mid,
		Group_name: site.Group_name,
		Group_logo: site.Group_logo,
		Group_extras: entities.GroupExtras{
			Open:      site.Open,
			Close:     site.Close,
			Address:   site.Address,
			Detail:    site.Detail,
			Estimate:  site.Estimate,
			Latitude:  site.Lat,
			Longitude: site.Long,
		},
		Trf: entities.SiteTrf{
			Adult: dataAdult,
			Child: dataChild,
		},
	}

	return &resp, "01", "Success get detail site data", true
}

//SelectTrip : select data trip
func SelectTrip(token *entities.Users, page int, size int) (*[]entities.TrxList, string, string, bool, int, int, int) {

	var trip []entities.TripTrxModel
	var counTrip []entities.TripModel

	offset := (page - 1) * size
	limit := size

	if err := db.DB[0].Select(`tp_id`).Where("tp_agent_id = ?", token.Typeid).Find(&counTrip).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Data transaction not found (" + err.Error() + ")", false, 0, 0, 0
	}

	if err := db.DB[0].Select(`tp_id, tp_status, tp_invoice, tp_number, tp_duration, tp_total_amount, trip_planner.created_at,
								agent_name,
								COALESCE(tp_start_date::text, '') as tp_start_date, tp_agent_id,
								COALESCE(tp_end_date::text, '') as tp_end_date,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'title' as text), '') as title,
								COALESCE(cast(tp_contact ->>'full_name' as text), '') as full_name,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'phone' as text), '') as phone,
								COALESCE(cast(tp_contact ->>'address' as text), '') as address`).Where("tp_agent_id = ?", token.Typeid).Joins("inner join master_agents on agent_id = tp_agent_id").Order("tp_invoice desc").Limit(limit).Offset(offset).Find(&trip).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Data transaction not found (" + err.Error() + ")", false, 0, 0, 0
	}

	var resp []entities.TrxList
	var status string

	for _, data := range trip {

		createdat, _ := time.Parse(time.RFC3339, data.Created_at)
		now := time.Now()

		if data.Tp_status == 2 && now.After(createdat) {
			data.Expired = true
		} else {
			data.Expired = false
		}

		if data.Tp_status == 1 {
			status = "Transaction Saved"
		} else if data.Tp_status == 2 && data.Expired == false {
			status = "Payment Pending"
		} else if data.Tp_status == 3 {
			status = "Payment Success"
		} else if data.Tp_status == 2 && data.Expired == true {
			status = "Expired at " + createdat.Format("2006-01-02")
		} else {
			status = "Unknown"
		}

		var persons []entities.TripPersonTrxModel

		if err := db.DB[0].Select(`tpp_id, tpp_name,
					CASE 
						WHEN tpp_type = 1 THEN 'adult'
						WHEN tpp_type = 2 THEN 'child'
						else 'unknown'
					end as "type",
					tpp_qr,
					COALESCE(cast(tpp_extras ->>'id' as text), '') as id_number,
					COALESCE(cast(tpp_extras ->>'title' as text), '') as title,
					COALESCE(cast(tpp_extras ->>'typeid' as text), '') as type_id`).Where("tpp_tp_id = ?", data.Tp_id).Find(&persons).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "04", "Data person not found (" + err.Error() + ")", false, 0, 0, 0
		}

		var respPerson []entities.TrxPerson

		for _, person := range persons {
			var dests []entities.TripDestinationTrxModel

			if err := db.DB[0].Select(`tpd_group_mid,
										trf_name, group_name,
										tpd_amount, tpd_date::text, 
										tpd_exp_date::text, 
										tpd_duration`).Where("tpd_tpp_id = ?", person.Tpp_id).Joins("inner join master_tariff on trf_id = tpd_trf_id").Joins("inner join master_group on group_mid = tpd_group_mid").Find(&dests).Error; gorm.IsRecordNotFoundError(err) {
				return nil, "05", "Data destination not found (" + err.Error() + ")", false, 0, 0, 0
			}

			tmpPerson := entities.TrxPerson{
				Id_number:   person.Id_number,
				Title:       person.Title,
				Tpp_name:    person.Tpp_name,
				Tpp_qr:      "TRP" + person.Tpp_qr,
				Type:        person.Type,
				Type_id:     person.Type_id,
				Destination: dests,
			}

			respPerson = append(respPerson, tmpPerson)
		}

		var grups []entities.TripGrupName

		if err := db.DB[0].Select(`distinct group_name`).Where("tp_id = ?", data.Tp_id).Joins("inner join trip_planner_person on tp_id = tpp_tp_id").Joins("inner join trip_planner_destination on tpp_id = tpd_tpp_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("group_name").Find(&grups).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "06", "Data group not found (" + err.Error() + ")", false, 0, 0, 0
		}

		var dest string

		for _, grup := range grups {
			dest += grup.Group_name + ", "
		}

		r := []rune(dest)

		if len(r) > 2 {
			dest = string(r[:len(r)-2])
		}

		tmpResp := entities.TrxList{
			Tp_number:       data.Tp_number,
			Tp_invoice:      data.Tp_invoice,
			Tp_duration:     data.Tp_duration,
			Tp_start_date:   data.Tp_start_date,
			Tp_end_date:     data.Tp_end_date,
			Tp_status:       data.Tp_status,
			Expired:         data.Expired,
			Created_at:      data.Created_at,
			Tp_id:           data.Tp_id,
			Status_name:     status,
			Tp_total_amount: data.Tp_total_amount,
			Tp_agent_id:     data.Tp_agent_id,
			Agent_name:      data.Agent_name,
			Destination:     dest,
			Contact: &entities.TrxContact{
				Email:    data.Email,
				Address:  data.Address,
				Fullname: data.Fullname,
				Phone:    data.Phone,
				Title:    data.Title,
			},
			Person: respPerson,
		}

		resp = append(resp, tmpResp)
	}

	totalData := len(counTrip)
	currentData := len(trip)
	rawPages := float64(totalData) / float64(size)
	totalPages := int(math.Ceil(rawPages))

	return &resp, "01", "Success get trx data", true, totalData, totalPages, currentData
}
