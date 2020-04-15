package repositories

import (
	"encoding/json"
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"

	uuid "github.com/satori/go.uuid"
)

// CheckoutB2B : b2b checkout ticket
func CheckoutB2B(token *entities.Users, r *entities.CheckOutReq) (map[string]interface{}, string, string, bool) {
	rTrip, err := json.Marshal(&r.Trip)
	if err != nil {
		return nil, "99", "Failed to parse json key trip (" + err.Error() + ")", false
	}

	rPerson, err := json.Marshal(&r.Person)
	if err != nil {
		return nil, "99", "Failed to parse json key person (" + err.Error() + ")", false
	}

	rContact, err := json.Marshal(&r.Header.Contact)
	if err != nil {
		return nil, "99", "Failed to parse json key contact (" + err.Error() + ")", false
	}

	var jTrip, jPerson, jContact string
	var dataQR []entities.QRTripRes
	jTrip = string(rTrip)
	jPerson = string(rPerson)

	jContact = string(rContact)
	jExtras := `{"trip":` + jTrip + `,` + `"person":` + jPerson + `}`

	tpID := uuid.NewV4()
	stan := int(time.Now().Unix())
	tpInvoice, err := strconv.Atoi(r.Header.Order)
	if err != nil {
		return nil, "99", "Failed to parse invoice order (" + err.Error() + ")", false
	}

	tripPlanner := entities.TripModel{
		Tp_id:           tpID,
		Tp_adult:        len(r.Person.Adult),
		Tp_child:        len(r.Person.Child),
		Tp_contact:      jContact,
		Tp_duration:     r.Header.Duration,
		Tp_start_date:   r.Header.StartDate,
		Tp_end_date:     r.Header.EndDate,
		Tp_invoice:      tpInvoice,
		Tp_number:       r.Header.InvNumber,
		Tp_src_type:     5,
		Tp_stan:         stan,
		Tp_status:       2,
		Tp_total_amount: r.Header.TotalAmount,
		Tp_extras:       jExtras,
		Tp_user_id:      token.ID,
		Created_at:      time.Now().Format("2006-01-02 15:04:05"),
	}

	db.DB[0].NewRecord(tripPlanner)

	if err := db.DB[0].Create(&tripPlanner).Error; err != nil {
		return nil, "02", "Error when inserting trip planner data (" + err.Error() + ")", false
	}

	for _, person := range r.Person.Adult {
		tppID := uuid.NewV4()
		qrCode := tppID.String() + `#` + strconv.Itoa(stan)
		tripPerson := entities.PersonModel{
			Tpp_id:     tppID,
			Tpp_tp_id:  tpID,
			Tpp_name:   person.Name,
			Tpp_type:   1,
			Tpp_qr:     qrCode,
			Tpp_extras: `{"id": "` + person.ID + `", "title": "` + person.Title + `", "typeid": "` + person.TypeID + `"}`,
			Created_at: time.Now().Format("2006-01-02 15:04:05"),
		}

		db.DB[0].NewRecord(tripPerson)

		if err := db.DB[0].Create(&tripPerson).Error; err != nil {
			return nil, "03", "Error when inserting trip planner person data (" + err.Error() + ")", false
		}

		var dataTrip []entities.TripDay

		for _, trip := range r.Trip {

			var dataDest []entities.Dest

			for _, dest := range trip.Destination {
				tpdID := uuid.NewV4()
				tripDest := entities.DestinationModel{
					Tpd_id:        tpdID,
					Tpd_tpp_id:    tppID,
					Tpd_amount:    dest.TrfAdult,
					Tpd_date:      trip.Date,
					Tpd_exp_date:  trip.Date + " 23:59:59",
					Tpd_day:       trip.Day,
					Tpd_duration:  dest.Duration,
					Tpd_trf_id:    dest.Trf_id_adult,
					Tpd_group_mid: dest.Mid,
					Tpd_extras:    `{"operational": "` + dest.Operational + `"}`,
					Created_at:    time.Now().Format("2006-01-02 15:04:05"),
				}

				db.DB[0].NewRecord(tripDest)

				if err := db.DB[0].Create(&tripDest).Error; err != nil {
					return nil, "04", "Error when inserting trip planner destination adult data (" + err.Error() + ")", false
				}

				tmpDest := entities.Dest{
					GroupName:   dest.GroupName,
					Duration:    dest.Duration,
					Mid:         dest.Mid,
					Operational: dest.Operational,
					Amount:      dest.TrfAdult,
					Trf_id:      dest.Trf_id_adult,
				}

				dataDest = append(dataDest, tmpDest)
			}

			tmpTrip := entities.TripDay{
				Day:         trip.Day,
				Date:        trip.Date,
				Tanggal:     trip.Tanggal,
				ExpiredDate: trip.Date + " 23:59:59",
				Destination: dataDest,
			}

			dataTrip = append(dataTrip, tmpTrip)
		}

		tmpPerson := entities.QRTripRes{
			ID:     person.ID,
			Name:   person.Name,
			QRcode: "TRP" + qrCode,
			Title:  person.Title,
			Type:   "adult",
			TypeID: person.TypeID,
			Trip:   dataTrip,
		}

		dataQR = append(dataQR, tmpPerson)
	}

	for _, child := range r.Person.Child {
		tppID := uuid.NewV4()
		qrCode := tppID.String() + `#` + strconv.Itoa(stan)
		tripPerson := entities.PersonModel{
			Tpp_id:     tppID,
			Tpp_tp_id:  tpID,
			Tpp_name:   child.Name,
			Tpp_type:   2,
			Tpp_qr:     qrCode,
			Tpp_extras: `{"id": "` + child.ID + `", "title": "` + child.Title + `", "typeid": "` + child.TypeID + `"}`,
			Created_at: time.Now().Format("2006-01-02 15:04:05"),
		}

		db.DB[0].NewRecord(tripPerson)

		if err := db.DB[0].Create(&tripPerson).Error; err != nil {
			return nil, "03", "Error when inserting trip planner child data (" + err.Error() + ")", false
		}

		var dataTrip []entities.TripDay

		for _, trip := range r.Trip {

			var dataDest []entities.Dest

			for _, dest := range trip.Destination {
				tpdID := uuid.NewV4()
				tripDest := entities.DestinationModel{
					Tpd_id:        tpdID,
					Tpd_tpp_id:    tppID,
					Tpd_amount:    dest.TrfChild,
					Tpd_date:      trip.Date,
					Tpd_exp_date:  trip.Date + " 23:59:59",
					Tpd_day:       trip.Day,
					Tpd_duration:  dest.Duration,
					Tpd_trf_id:    dest.Trf_id_child,
					Tpd_group_mid: dest.Mid,
					Tpd_extras:    `{"operational": "` + dest.Operational + `"}`,
					Created_at:    time.Now().Format("2006-01-02 15:04:05"),
				}

				db.DB[0].NewRecord(tripDest)

				if err := db.DB[0].Create(&tripDest).Error; err != nil {
					return nil, "04", "Error when inserting trip planner destination adult data (" + err.Error() + ")", false
				}

				tmpDest := entities.Dest{
					GroupName:   dest.GroupName,
					Duration:    dest.Duration,
					Mid:         dest.Mid,
					Operational: dest.Operational,
					Amount:      dest.TrfChild,
					Trf_id:      dest.Trf_id_child,
				}

				dataDest = append(dataDest, tmpDest)
			}

			tmpTrip := entities.TripDay{
				Day:         trip.Day,
				Date:        trip.Date,
				Tanggal:     trip.Tanggal,
				ExpiredDate: trip.Date + " 23:59:59",
				Destination: dataDest,
			}

			dataTrip = append(dataTrip, tmpTrip)
		}

		tmpPerson := entities.QRTripRes{
			ID:     child.ID,
			Name:   child.Name,
			QRcode: "TRP" + qrCode,
			Title:  child.Title,
			Type:   "child",
			TypeID: child.TypeID,
			Trip:   dataTrip,
		}

		dataQR = append(dataQR, tmpPerson)
	}

	return map[string]interface{}{
		"tp_id":  tpID,
		"person": dataQR,
	}, "01", "Checkout success", true
}
