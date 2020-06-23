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
func SelectCluster(token *entities.Users) (*[]entities.Cluster, string, string, bool) {

	var cluster []entities.GrupModel

	if err := db.DB[1].Where("depth = 2 AND deleted_at is NULL").Find(&cluster).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Cluster not found (" + err.Error() + ")", false
	}

	var resp []entities.Cluster

	for _, data := range cluster {
		var sites []entities.GrupModel

		if err := db.DB[1].Select("*, COALESCE(cast(group_extras ->>'estimate' as text), '') as group_estimate").Where("depth = 3 AND parent_id = ? AND deleted_at is NULL", data.Group_id).Find(&sites).Error; err != nil {
			return nil, "03", "Error when fetching site data (" + err.Error() + ")", false
		}

		var siteResp []entities.Site

		for _, site := range sites {
			tmpSite := entities.Site{
				SiteID:        site.Group_id,
				SiteMID:       site.Group_mid,
				SiteName:      site.Group_name,
				SiteLogo:      site.Group_logo,
				SiteEstimated: site.Group_estimate,
			}

			siteResp = append(siteResp, tmpSite)
		}

		tmpResp := entities.Cluster{
			ClusterID:   data.Group_id,
			ClusterMID:  data.Group_mid,
			ClusterName: data.Group_name,
			ClusterLogo: data.Group_logo,
			Site:        siteResp,
		}

		resp = append(resp, tmpResp)
	}

	return &resp, "01", "Success get cluster data", true
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

	if err := db.DB[0].Select(`tp_id, tp_status, tp_invoice, tp_number, tp_start_date, tp_end_date, tp_duration, tp_total_amount,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'title' as text), '') as title,
								COALESCE(cast(tp_contact ->>'fullname' as text), '') as fullname,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'phone' as text), '') as phone,
								COALESCE(cast(tp_contact ->>'address' as text), '') as address`).Where("tp_agent_id = ?", token.Typeid).Order("tp_invoice desc").Limit(limit).Offset(offset).Find(&trip).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Data transaction not found (" + err.Error() + ")", false, 0, 0, 0
	}

	var resp []entities.TrxList
	var status string

	for _, data := range trip {
		if data.Tp_status == 1 {
			status = "bucket"
		} else if data.Tp_status == 2 {
			status = "checkout"
		} else if data.Tp_status == 3 {
			status = "purchased"
		} else {
			status = "unknown"
		}

		var grups []entities.TripGrupName

		if err := db.DB[0].Select(`distinct group_name`).Where("tp_id = ?", data.Tp_id).Joins("inner join trip_planner_person on tp_id = tpp_tp_id").Joins("inner join trip_planner_destination on tpp_id = tpd_tpp_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("group_name").Find(&grups).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "03", "Data group not found (" + err.Error() + ")", false, 0, 0, 0
		}

		var dest string

		for _, grup := range grups {
			dest += grup.Group_name + ", "
		}

		r := []rune(dest)
		dest = string(r[:len(r)-2])

		tmpResp := entities.TrxList{
			Tp_number:       data.Tp_number,
			Tp_invoice:      data.Tp_invoice,
			Tp_duration:     data.Tp_duration,
			Tp_start_date:   data.Tp_start_date,
			Tp_end_date:     data.Tp_end_date,
			Tp_status:       status,
			Tp_total_amount: data.Tp_total_amount,
			Destination:     dest,
			Contact: &entities.TrxContact{
				Email:    data.Email,
				Address:  data.Address,
				Fullname: data.Fullname,
				Phone:    data.Phone,
				Title:    data.Title,
			},
		}

		resp = append(resp, tmpResp)
	}

	totalData := len(counTrip)
	currentData := len(trip)
	rawPages := float64(totalData) / float64(size)
	totalPages := int(math.Ceil(rawPages))

	return &resp, "01", "Success get trx data", true, totalData, totalPages, currentData
}
