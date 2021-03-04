package repositories

import (
	"fmt"
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// InsertTrx : insert to table trx
func InsertTrx(token *entities.Users, r *requests.TrxReq) (*requests.TrxResp, string, string, bool) {
	if len(r.Customer) == 0 {
		return nil, "99", "Customer is required", false
	}

	// if r.TotalAmount == 0 {
	// 	return nil, "99", "Total amount is required", false
	// }
	if r.DestQty == 0 {
		return nil, "99", "Destination qty is required", false
	}

	if r.SourceType == 0 {
		return nil, "99", "Source type is required", false
	}

	if r.StartDate == "" {
		return nil, "99", "Start date is required", false
	}

	if r.EndDate == "" {
		return nil, "99", "End date is required", false
	}

	if r.TransactionStatus != 1 || r.TransactionStatus != 2 {
		return nil, "99", "Invalid transaction_value status", false
	}

	// var vis []requests.TrxVisit
	// var totPay float32
	var name, phone, email string
	var totalAmount float32

	tpID := uuid.NewV4()
	stan := int(time.Now().Unix())
	bNumber := "MB2B." + string(time.Now().Format("020106")) + "." + strconv.Itoa(stan)
	tripPlanner := entities.TrpTrxModel{
		Tp_id:           tpID,
		Tp_contact:      "{}",
		Tp_duration:     (helper.DaysBetween(helper.Date(r.StartDate), helper.Date(r.EndDate))) + 1,
		Tp_start_date:   r.StartDate,
		Tp_end_date:     r.EndDate,
		Tp_number:       bNumber,
		Tp_src_type:     r.SourceType,
		Tp_stan:         stan,
		Tp_status:       r.TransactionStatus,
		Tp_total_amount: 0,
		Tp_user_id:      token.ID,
		Tp_agent_id:     token.Typeid,
		Created_at:      time.Now().Format("2006-01-02 15:04:05"),
	}

	db.DB[0].NewRecord(tripPlanner)

	if err := db.DB[0].Create(&tripPlanner).Error; err != nil {
		return nil, "02", "Error when inserting trip planner data (" + err.Error() + ")", false
	}

	for _, cust := range r.Customer {
		if cust.IsPic == true {
			var trp entities.TrpTrxModel

			if err := db.DB[0].Model(&trp).Where("tp_id = ?", tpID).UpdateColumn("tp_contact", `{"nationality":"`+cust.Nationality+`", "region":"`+cust.Region+
				`", "typeid":"`+cust.IDType+`", "id":"`+cust.IDNumber+
				`", "idname":"`+cust.Name+`", "type":"`+cust.Type+
				`", "title":"`+cust.Title+`", "email":"`+cust.Email+
				`", "phone":"`+cust.Phone+`", "pic":`+strconv.FormatBool(cust.IsPic)+`}`).Error; err != nil {
				return nil, "03", "Error when updating trip data (" + err.Error() + ")", false
			}

			name = cust.Name
			phone = cust.Phone
			email = cust.Email
		}
		tppID := uuid.NewV4()
		qrCode := tppID.String() + `#` + strconv.Itoa(stan)
		var custType int

		if cust.Type == "adult" {
			custType = 1
		} else {
			custType = 2
		}

		tripPerson := entities.PersonModel{
			Tpp_id:    tppID,
			Tpp_tp_id: tpID,
			Tpp_name:  cust.Name,
			Tpp_type:  custType,
			Tpp_qr:    qrCode,
			Tpp_extras: `{"nationality":"` + cust.Nationality + `", "region":"` + cust.Region +
				`", "typeid":"` + cust.IDType + `", "id":"` + cust.IDNumber +
				`", "idname":"` + cust.Name + `", "type":"` + cust.Type +
				`", "title":"` + cust.Title + `", "email":"` + cust.Email +
				`", "phone":"` + cust.Phone + `", "pic":` + strconv.FormatBool(cust.IsPic) + `}`,
			Created_at: time.Now().Format("2006-01-02 15:04:05"),
		}

		db.DB[0].NewRecord(tripPerson)

		if err := db.DB[0].Create(&tripPerson).Error; err != nil {
			return nil, "03", "Error when inserting trip planner person data (" + err.Error() + ")", false
		}

		for _, trip := range cust.Trip {
			for _, dest := range trip.Ticket {
				var getExp entities.GetExp
				var duration int

				if err := db.DB[1].Select(`cast(trf_condition ->> 'expiredQr' as int) as expired,
											group_mid, cast(group_extras ->> 'estimate' as text) duration`).Joins("join master_group on group_id = trf_group_id").Where("trf_id = ?", dest.TrfID).Find(&getExp).Error; err != nil {
					return nil, "04", "Couldn't find tariff with id: " + strconv.Itoa(dest.TrfID) + " (" + err.Error() + ")", false
				}

				if getExp.Duration == "" {
					duration = 1
				} else {
					durconv, err := strconv.Atoi(getExp.Duration)

					if err != nil {
						return nil, "05", err.Error(), false
					}

					duration = durconv
				}

				// r.Trip[len(r.Trip)-1].TripDate
				t, _ := time.Parse("2006-01-02", r.EndDate)

				// dayExp := t.Add(time.Hour*time.Duration((getExp.Expired*24)-24)).Format("2006-01-02") + " 23:59:59"

				dayExp := t.Format("2006-01-02") + " 23:59:59"

				ext := `{"original_amount":` + fmt.Sprintf("%g", dest.Amount) + `, "discount":` + fmt.Sprintf("%g", dest.Discount) + `}`

				tpdID := uuid.NewV4()
				tripDes := entities.DestinationModel{
					Tpd_id:        tpdID,
					Tpd_tpp_id:    tppID,
					Tpd_amount:    dest.NettAmount,
					Tpd_date:      trip.TripDate,
					Tpd_day:       trip.TripDay,
					Tpd_duration:  duration,
					Tpd_exp_date:  dayExp,
					Tpd_group_mid: getExp.Group_mid,
					Tpd_trf_id:    dest.TrfID,
					Tpd_extras:    ext,
					Created_at:    time.Now().Format("2006-01-02 15:04:05"),
				}

				db.DB[0].NewRecord(tripDes)

				if err := db.DB[0].Create(&tripDes).Error; err != nil {
					return nil, "06", "Error when inserting trip planner destination data (" + err.Error() + ")", false
				}

				totalAmount += dest.NettAmount
			}

			// tmpVisit := requests.TrxVisit{
			// 	TripDate:    trip.TripDate,
			// 	TripDay:     trip.TripDay,
			// 	Ticket:      trip.Ticket,
			// }

			// vis = append(vis, tmpVisit)
		}

	}

	if err := db.DB[0].Model(&tripPlanner).Where("tp_id = ?", tpID).Update("tp_total_amount", totalAmount).Error; err != nil {
		return nil, "03", "Error when updating trip data tp_total_amount (" + err.Error() + ")", false
	}

	res := requests.TrxResp{
		BookingNumber: bNumber,
		PayTotal:      totalAmount,
		Name:          name,
		Email:         email,
		Phone:         phone,
		Cust:          r.Customer,
	}

	return &res, "01", "Transaction success", true
}

// UpdateTrx : update table trx
func UpdateTrx(token *entities.Users, r *requests.TrxReqUpdate) (*[]entities.UpdateTrxResp, string, string, bool) {
	if len(r.Trx) == 0 {
		return nil, "99", "Trip is required", false
	}

	if r.Status == 0 {
		return nil, "99", "Status is required", false
	}

	var resp []entities.UpdateTrxResp

	for _, trx := range r.Trx {
		update, msg := updateTableTrx(trx.BookingNumber, r.Status, trx.PaymentMethod)
		if update == false {
			return nil, "02", "Failed to update data (" + msg + ")", false
		}

		var trp entities.TrpTrxModel

		if err := db.DB[0].Select(`tp_status, tp_number`).Where("tp_number = ?", trx.BookingNumber).Find(&trp).Error; err != nil {
			return nil, "03", "Couldn't find invoice with number: " + trx.BookingNumber + " (" + err.Error() + ")", false
		}

		var status string

		if trp.Tp_status == 1 {
			status = "Draft"
		} else if trp.Tp_status == 2 {
			status = "Purchase"
		} else if trp.Tp_status == 3 {
			status = "Paid"
		} else {
			status = "Unknown"
		}

		tmpResp := entities.UpdateTrxResp{
			InvNumber: trp.Tp_number,
			Status:    status,
		}

		resp = append(resp, tmpResp)
	}

	return &resp, "01", "Transaction updated successfully", true
}

// UpdateTrxPayment : update table trx
func UpdateTrxPayment(token *entities.Users, r *requests.TrxReqUpdate) (*[]entities.UpdateTrxResp, string, string, bool) {
	if len(r.Trx) == 0 {
		return nil, "99", "Transactions is required", false
	}

	var resp []entities.UpdateTrxResp

	for _, trx := range r.Trx {
		update, msg := updateTableTrx(trx.BookingNumber, 3, trx.PaymentMethod)
		if update == false {
			return nil, "02", "Failed to update data (" + msg + ")", false
		}

		var trp entities.TrpTrxModel

		if err := db.DB[0].Select(`tp_status, tp_number`).Where("tp_number = ?", trx.BookingNumber).Find(&trp).Error; err != nil {
			return nil, "03", "Couldn't find invoice with number: " + trx.BookingNumber + " (" + err.Error() + ")", false
		}

		var status string

		if trp.Tp_status == 1 {
			status = "Draft"
		} else if trp.Tp_status == 2 {
			status = "Purchase"
		} else if trp.Tp_status == 3 {
			status = "Paid"
		} else {
			status = "Unknown"
		}

		tmpResp := entities.UpdateTrxResp{
			InvNumber: trp.Tp_number,
			Status:    status,
		}

		resp = append(resp, tmpResp)
	}

	return &resp, "01", "Payment succeed", true
}

func updateTableTrx(inv string, status int, paymentMethod string) (bool, string) {
	var trp entities.TrpTrxModel

	if paymentMethod == "" {
		if err := db.DB[0].Model(&trp).Where("tp_number = ?", inv).Updates(map[string]interface{}{"tp_status": status, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
			return false, err.Error()
		}
	} else {
		if err := db.DB[0].Model(&trp).Where("tp_number = ?", inv).Updates(map[string]interface{}{"tp_status": status, "tp_payment_method": paymentMethod, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
			return false, err.Error()
		}
	}

	return true, ""
}

//GetQR : select data trip
func GetQR(token *entities.Users, r *requests.TrxQReq) (*entities.TrxList, string, string, bool) {
	if r.Inv == "" {
		return nil, "99", "Invoice number is required", false
	}

	var trip entities.TripTrxModel

	if err := db.DB[0].Select(`tp_id, tp_status, tp_invoice, tp_number, tp_duration, tp_total_amount,
								agent_name,
								COALESCE(tp_start_date::text, '') as tp_start_date, tp_agent_id,
								COALESCE(tp_end_date::text, '') as tp_end_date,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'title' as text), '') as title,
								coalesce(cast(tp_contact ->>'idname' as text), '') as fullname,
								COALESCE(cast(tp_contact ->>'email' as text), '') as email,
								COALESCE(cast(tp_contact ->>'phone' as text), '') as phone,
								COALESCE(cast(tp_contact ->>'region' as text), '') as address`).Where("tp_agent_id = ? and tp_number = ?", token.Typeid, r.Inv).Joins("inner join master_agents on agent_id = tp_agent_id").Find(&trip).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Data transaction not found (" + err.Error() + ")", false
	}

	var status string

	if trip.Tp_status == 1 {
		status = "Draft"
	} else if trip.Tp_status == 2 {
		status = "Purchase"
	} else if trip.Tp_status == 3 {
		status = "Paid"
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
					COALESCE(cast(tpp_extras ->>'typeid' as text), '') as type_id`).Where("tpp_tp_id = ?", trip.Tp_id).Find(&persons).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "04", "Data person not found (" + err.Error() + ")", false
	}

	var respPerson []entities.TrxPerson

	for _, person := range persons {
		var dests []entities.TripDestinationTrxModel

		if err := db.DB[0].Select(`tpd_group_mid,
										trf_name, group_name,
										tpd_amount, tpd_date::text, 
										tpd_exp_date::text, 
										tpd_duration`).Where("tpd_tpp_id = ?", person.Tpp_id).Joins("inner join master_tariff on trf_id = tpd_trf_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("tpd_date").Find(&dests).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "05", "Data destination not found (" + err.Error() + ")", false
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

	if err := db.DB[0].Select(`distinct group_name`).Where("tp_id = ?", trip.Tp_id).Joins("inner join trip_planner_person on tp_id = tpp_tp_id").Joins("inner join trip_planner_destination on tpp_id = tpd_tpp_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("group_name").Find(&grups).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "06", "Data group not found (" + err.Error() + ")", false
	}

	var dest string

	for _, grup := range grups {
		dest += grup.Group_name + ", "
	}

	z := []rune(dest)

	if len(z) > 2 {
		dest = string(z[:len(z)-2])
	}

	resp := entities.TrxList{
		Tp_number:       trip.Tp_number,
		Tp_invoice:      trip.Tp_invoice,
		Tp_duration:     trip.Tp_duration,
		Tp_start_date:   trip.Tp_start_date,
		Tp_end_date:     trip.Tp_end_date,
		Tp_status:       trip.Tp_status,
		Tp_id:           trip.Tp_id,
		Status_name:     status,
		Tp_total_amount: trip.Tp_total_amount,
		Tp_agent_id:     trip.Tp_agent_id,
		Agent_name:      trip.Agent_name,
		Destination:     dest,
		Contact: &entities.TrxContact{
			Email:    trip.Email,
			Address:  trip.Address,
			Fullname: trip.Fullname,
			Phone:    trip.Phone,
			Title:    trip.Title,
		},
		Person: respPerson,
	}

	return &resp, "01", "Success get trx data", true
}

//GetTrxByNumber : select data trip by number
func GetTrxByNumber(token *entities.Users, r *requests.TrxQReq) (*entities.RespTrxNum, string, string, bool) {
	if r.TrxNum == "" {
		return nil, "99", "Trx number is required", false
	}

	var trip entities.TripTrxModel

	if err := db.DB[0].Select(`tp_id, tp_status, tp_number, tp_duration, tp_total_amount,
								agent_name,
								COALESCE(tp_start_date::text, '') as tp_start_date,
								COALESCE(tp_end_date::text, '') as tp_end_date,
								COALESCE(trip_planner.created_at::text, '') as created_at,
								coalesce(cast(tp_contact ->>'idname' as text), '') as fullname`).Where("tp_number = ?", r.TrxNum).Joins("inner join master_agents on agent_id = tp_agent_id").Find(&trip).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Data transaction not found (" + err.Error() + ")", false
	}

	var status string

	if trip.Tp_status == 1 {
		status = "Created"
	} else if trip.Tp_status == 2 {
		status = "Purchase"
	} else if trip.Tp_status == 3 {
		status = "Paid"
	} else {
		status = "Unknown"
	}

	var persons []entities.TripPersonTrxModel

	if err := db.DB[0].Select(`tpp_id, tpp_name,
					CASE 
						WHEN tpp_type = 1 THEN 'Adult'
						WHEN tpp_type = 2 THEN 'Child'
						else 'unknown'
					end as "type",
					tpp_qr`).Where("tpp_tp_id = ?", trip.Tp_id).Find(&persons).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "04", "Data person not found (" + err.Error() + ")", false
	}

	var respPerson []entities.RespPersonNum

	for _, person := range persons {
		var dests []entities.DestinationTrxModel

		if err := db.DB[0].Select(`trf_name, group_name,
										tpd_amount, 
										coalesce(cast(tpd_extras ->>'original_amount' as float), 0) as bruto,
										coalesce(cast(tpd_extras ->>'discount' as float), 0) as disc`).Where("tpd_tpp_id = ?", person.Tpp_id).Joins("inner join master_tariff on trf_id = tpd_trf_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("tpd_date").Find(&dests).Error; gorm.IsRecordNotFoundError(err) {
			return nil, "05", "Data destination not found (" + err.Error() + ")", false
		}

		tmpPerson := entities.RespPersonNum{
			Tpp_name:    person.Tpp_name,
			Tpp_qr:      "TRP" + person.Tpp_qr,
			Type:        person.Type,
			Destination: dests,
		}

		respPerson = append(respPerson, tmpPerson)
	}

	// var grups []entities.TripGrupName

	// if err := db.DB[0].Select(`distinct group_name`).Where("tp_id = ?", trip.Tp_id).Joins("inner join trip_planner_person on tp_id = tpp_tp_id").Joins("inner join trip_planner_destination on tpp_id = tpd_tpp_id").Joins("inner join master_group on group_mid = tpd_group_mid").Order("group_name").Find(&grups).Error; gorm.IsRecordNotFoundError(err) {
	// 	return nil, "06", "Data group not found (" + err.Error() + ")", false
	// }

	// var dest string

	// for _, grup := range grups {
	// 	dest += grup.Group_name + ", "
	// }

	// z := []rune(dest)

	// if len(z) > 2 {
	// 	dest = string(z[:len(z)-2])
	// }

	resp := entities.RespTrxNum{
		Tp_number:       trip.Tp_number,
		Tp_duration:     trip.Tp_duration,
		Tp_start_date:   trip.Tp_start_date,
		Tp_end_date:     trip.Tp_end_date,
		Tp_status:       status,
		Agent:           trip.Agent_name,
		Fullname:        trip.Fullname,
		Tp_trx_date:     trip.Created_at,
		Tp_total_amount: trip.Tp_total_amount,
		Person:          respPerson,
	}

	return &resp, "01", "Success get trx data", true
}
