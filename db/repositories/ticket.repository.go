package repositories

import (
	"strings"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
)

var qrPrefix = "AINO"
var prefixLen = len(qrPrefix)
var weekday = strings.ToUpper("%" + time.Now().Format("Mon") + "%")

// GetTicket : get ticket list for device
func GetTicket(r interface{}, token *entities.Users) (*[]entities.MasterTicket, string, string, bool) {
	var masterTicket []entities.MasterTicket
	var mTickets []entities.MasterTicket
	var tickets []entities.Ticket

	mbmid := r.(map[string]interface{})["mbmid"]

	if mbmid == nil {
		return nil, "99", "Mbmid are required", false
	}

	query := `SELECT 
				COALESCE(cast(a.trf_id as varchar(255)), '') as trf_id,
				COALESCE(a.trf_code, '') as trf_code,
				COALESCE(cast(a.trf_group_id as varchar(255)), '') as trf_group_id,
				COALESCE(f.trfftype_name, '') as trf_trftype,
				COALESCE(a.trf_name, '') as trf_name,
				COALESCE(a.trf_value::text, '') as trf_value,
				COALESCE(a.trf_start_date::text, '') as trf_start_date,
				COALESCE(a.trf_end_date::text, '') as trf_end_date,
				COALESCE(a.trf_priority::text, '') as trf_priority,
				COALESCE(a.trf_release::text, '') as trf_release,
				COALESCE(a.trf_currency_code, '') as trf_currency_code,
				COALESCE(a.trf_qty, '1') as trf_qty,
				COALESCE(cast(trf_condition->>'day' as text), '') as day,
				COALESCE(cast(trf_condition->>'beginTime' as text), '') as begin_time,
				COALESCE(cast(trf_condition->>'endTime' as text), '') as end_time,
				COALESCE(cast(trf_condition->>'cardType' as text), '') as card_type,
				COALESCE(cast(trf_condition->>'expiredQr' as text), '') as expired_qr
				from public.master_tariff a
				join public.master_group d on a.trf_group_id = d.group_id
				left join public.master_tariff_type f on a.trf_trfftype_id = f.trfftype_id
				left join public.master_tariff_has_machine i on a.trf_id = i.trf_id 
				where (d.group_mid = ?)
				and a.deleted_at is null
				and a.trf_agent_id = ?`

	query += ` order by a.trf_id DESC;`
	db.DB[0].Raw(query, mbmid, token.Typeid).Scan(&masterTicket)

	if len(masterTicket) == 0 {
		return nil, "02", "Fare not found", false
	}

	for _, item := range masterTicket {
		var dataTicket []entities.Ticket
		db.DB[0].Raw(`select COALESCE(c.mtick_name, '') as mtick_name,
								COALESCE(c.mtick_code, '') as mtick_code,
								COALESCE(g.tctype_name, '') as mtick_type,
								COALESCE(k.ctg_name, '') as mtick_cat
						from master_ticket c
							left join master_tariffdet b on c.mtick_id = b.trfdet_mtick_id
							left join master_tariff a on b.trfdet_trf_id = a.trf_id
							left join public.master_ticket_type g on c.mtick_type_id = g.tctype_id
							left join public.master_ticket_category k on c.mtick_ctg_id = k.ctg_id
						where a.trf_id = ?`,
			item.Trf_id).Scan(&tickets)

		if len(tickets) == 0 {
			return nil, "03", "Ticket not found", false
		}

		for _, ticket := range tickets {
			tmpTicket := entities.Ticket{
				Mtick_name: ticket.Mtick_name,
				Mtick_code: ticket.Mtick_code,
				Mtick_type: ticket.Mtick_type,
				Mtick_cat:  ticket.Mtick_cat,
			}
			dataTicket = append(dataTicket, tmpTicket)
		}

		dayArr := strings.Split(item.Day, ",")
		// cardArr := strings.Split(item.Card_type, ",")
		cardArr := item.Card_type

		mTicket := entities.MasterTicket{
			Trf_id:            item.Trf_id,
			Trf_name:          item.Trf_name,
			Trf_code:          item.Trf_code,
			Trf_group_id:      item.Trf_group_id,
			Trf_trftype:       item.Trf_trftype,
			Trf_value:         item.Trf_value,
			Trf_start_date:    item.Trf_start_date,
			Trf_end_date:      item.Trf_end_date,
			Trf_priority:      item.Trf_priority,
			Trf_release:       item.Trf_release,
			Trf_currency_code: item.Trf_currency_code,
			Trf_qty:           item.Trf_qty,
			Ticket:            dataTicket,
			Trf_condition: &entities.Condition{
				Day:        dayArr,
				Begin_time: item.Begin_time,
				End_time:   item.End_time,
				Card_type:  cardArr,
				Expired_qr: item.Expired_qr,
			},
		}

		mTickets = append(mTickets, mTicket)
	}

	return &mTickets, "01", "Success get ticket list", true
}
