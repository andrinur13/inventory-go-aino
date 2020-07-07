package repositories

import (
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	uuid "github.com/satori/go.uuid"
)

// InsertTrx : insert to table trx
func InsertTrx(token *entities.Users, r *requests.TrxReq) (*requests.TrxResp, string, string, bool) {
	if len(r.Trip) == 0 {
		return nil, "99", "Trip is required", false
	}

	if len(r.Customer) == 0 {
		return nil, "99", "Customer is required", false
	}

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

	var vis []requests.TrxVisit
	var totPay float32
	var name, phone, email string

	for _, trip := range r.Trip {
		totPay += trip.TotalAmount
		tpID := uuid.NewV4()
		stan := int(time.Now().Unix())
		bNumber := "WB2B." + string(time.Now().Format("020106")) + "." + strconv.Itoa(stan)
		tripPlanner := entities.TrpTrxModel{
			Tp_id:           tpID,
			Tp_contact:      "{}",
			Tp_duration:     (helper.DaysBetween(helper.Date(r.StartDate), helper.Date(r.EndDate))) + 1,
			Tp_start_date:   r.StartDate,
			Tp_end_date:     r.EndDate,
			Tp_number:       bNumber,
			Tp_src_type:     r.SourceType,
			Tp_stan:         stan,
			Tp_status:       1,
			Tp_total_amount: trip.TotalAmount,
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

				// db.DB[0].Where("tp_id = ?", tpID).Find(&trp)
				// trp.Tp_contact = `{"nationality":"` + cust.Nationality + `", "region":"` + cust.Region +
				// 	`", "idtype":"` + cust.IDType + `", "idnumber":"` + cust.IDNumber +
				// 	`", "idname":"` + cust.Name + `", "type":"` + cust.Type +
				// 	`", "title":"` + cust.Title + `", "email":"` + cust.Email +
				// 	`", "phone":"` + cust.Phone + `", "pic":` + strconv.FormatBool(cust.IsPic) + `}`
				// if err := db.DB[0].Save(&trp).Error; err != nil {
				// 	return nil, "03", "Error when updating trip data (" + err.Error() + ")", false
				// }
				if err := db.DB[0].Model(&trp).Where("tp_id = ?", tpID).Update("tp_contact", `{"nationality":"`+cust.Nationality+`", "region":"`+cust.Region+
					`", "idtype":"`+cust.IDType+`", "idnumber":"`+cust.IDNumber+
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
					`", "idtype":"` + cust.IDType + `", "idnumber":"` + cust.IDNumber +
					`", "idname":"` + cust.Name + `", "type":"` + cust.Type +
					`", "title":"` + cust.Title + `", "email":"` + cust.Email +
					`", "phone":"` + cust.Phone + `", "pic":` + strconv.FormatBool(cust.IsPic) + `}`,
				Created_at: time.Now().Format("2006-01-02 15:04:05"),
			}

			db.DB[0].NewRecord(tripPerson)

			if err := db.DB[0].Create(&tripPerson).Error; err != nil {
				return nil, "03", "Error when inserting trip planner person data (" + err.Error() + ")", false
			}

			for _, dest := range trip.Ticket {
				var getExp entities.GetExp

				if err := db.DB[0].Select(`cast(trf_condition ->> 'expiredQr' as int) as expired`).Where("trf_id = ?", dest.TrfID).Find(&getExp).Error; err != nil {
					return nil, "04", "Couldn't find tariff with id: " + strconv.Itoa(dest.TrfID) + " (" + err.Error() + ")", false
				}

				t, _ := time.Parse("2006-01-02", trip.TripDate)

				dayExp := t.Add(time.Hour*time.Duration((getExp.Expired*24)-24)).Format("2006-01-02") + " 23:59:59"

				tpdID := uuid.NewV4()
				tripDes := entities.DestinationModel{
					Tpd_id:        tpdID,
					Tpd_tpp_id:    tppID,
					Tpd_amount:    dest.Amount,
					Tpd_date:      trip.TripDate,
					Tpd_day:       trip.TripDay,
					Tpd_duration:  dest.SiteDuration,
					Tpd_exp_date:  dayExp,
					Tpd_group_mid: dest.Mmid,
					Tpd_trf_id:    dest.TrfID,
					Tpd_extras:    `{}`,
					Created_at:    time.Now().Format("2006-01-02 15:04:05"),
				}

				db.DB[0].NewRecord(tripDes)

				if err := db.DB[0].Create(&tripDes).Error; err != nil {
					return nil, "03", "Error when inserting trip planner destination data (" + err.Error() + ")", false
				}
			}

		}

		tmpVisit := requests.TrxVisit{
			BookingNumber: bNumber,
			TotalAmount:   trip.TotalAmount,
			TripDate:      trip.TripDate,
			TripDay:       trip.TripDay,
			Ticket:        trip.Ticket,
		}

		vis = append(vis, tmpVisit)
	}

	res := requests.TrxResp{
		PayTotal: totPay,
		Name:     name,
		Email:    email,
		Phone:    phone,
		Cust:     r.Customer,
		Visit:    vis,
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
