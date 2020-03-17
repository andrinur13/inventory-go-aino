package repositories

import (
	"twc-ota-api/config"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"

	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
)

// BookingTicket : booking ticket
func BookingTicket(token *entities.Users, r *requests.BookingReq) (map[string]interface{}, string, string, bool) {
	if r.BookingNumber == "" {
		return nil, "99", "Format error, booking number is required", false
	}

	if r.BookingDate == "" {
		return nil, "99", "Format error, booking date is required", false
	}

	if r.Mbmid == "" {
		return nil, "99", "Format error, mbmid is required", false
	}

	if r.PayAmt == 0 {
		return nil, "99", "Format error, amount is required", false
	}

	if r.Emoney == 0 {
		return nil, "99", "Format error, emoney is required", false
	}

	if r.PayMethod == "" {
		return nil, "99", "Format error, payment method is required", false
	}

	if r.Email == "" {
		return nil, "99", "Format error, customer's email is required", false
	}

	if len(r.Trf) == 0 {
		return nil, "99", "Format error, please complete tarif's payload", false
	}

	booking := entities.Booking{
		Agent_id:               token.Typeid,
		Booking_number:         r.BookingNumber,
		Booking_date:           r.BookingDate,
		Booking_mid:            r.Mbmid,
		Booking_amount:         r.PayAmt,
		Booking_emoney:         r.Emoney,
		Booking_total_payment:  r.PayAmt,
		Booking_payment_method: r.PayMethod,
		Customer_email:         r.Email,
		Customer_phone:         r.Phone,
		Customer_username:      r.Username,
		Customer_note:          r.Note,
	}

	db.DB[0].NewRecord(booking)

	if err := db.DB[0].Create(&booking).Error; err != nil {
		return nil, "02", "Error when inserting booking data (" + err.Error() + ")", false
	}

	for _, trf := range r.Trf {
		if trf.TrfID == 0 || trf.TrfAmount == 0 || trf.TrfType == "" || trf.TrfQty == 0 || trf.TrfTotal == 0 {
			db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
			db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

			return nil, "99", "Format error, please complete tarif's payload", false
		}
		var tarif entities.TariffModel
		if err := db.DB[0].Where("trf_id = ? AND deleted_at IS NULL", trf.TrfID).Find(&tarif).Error; gorm.IsRecordNotFoundError(err) {
			db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
			db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

			return nil, "03", "Specified fare on request not found (" + err.Error() + ")", false
		}

		for i := 0; i < trf.TrfQty; i++ {
			bookdet := entities.Bookingdet{
				Bookingdet_booking_id: int(booking.Booking_id),
				Bookingdet_trf_id:     trf.TrfID,
				Bookingdet_trftype:    trf.TrfType,
				Bookingdet_amount:     trf.TrfAmount,
				Bookingdet_qty:        1,
				Bookingdet_total:      trf.TrfAmount,
			}
			db.DB[0].NewRecord(bookdet)

			if err := db.DB[0].Create(&bookdet).Error; err != nil {
				db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
				db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

				return nil, "04", "Error when inserting bookingdet data (" + err.Error() + ")", false
			}

			var tarifDet []entities.TariffDetModel

			if err := db.DB[0].Where("trfdet_trf_id = ?", trf.TrfID).Find(&tarifDet).Error; gorm.IsRecordNotFoundError(err) {
				db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
				db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

				return nil, "05", "Fare not found on master_tariffdet (" + err.Error() + ")", false
			}

			for _, tickets := range tarifDet {
				var ticket entities.TicketModel
				if err := db.DB[0].Where("mtick_id = ? AND deleted_at IS NULL", tickets.Trfdet_mtick_id).Find(&ticket).Error; gorm.IsRecordNotFoundError(err) {
					db.DB[0].Where("bookinglist_bookingdet_id = ?", bookdet.Bookingdet_id).Delete(entities.Bookinglist{})
					db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
					db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

					return nil, "06", "Specified ticket on master_tariffdet not found (" + err.Error() + ")", false
				}

				var grup entities.GrupModel

				db.DB[0].Table("master_group").Select("group_mid").Where("mtick_id = ?", ticket.Mtick_id).Joins("left join master_ticket on mtick_group_id = group_id").Scan(&grup)

				booklist := entities.Bookinglist{
					Bookinglist_bookingdet_id: int(bookdet.Bookingdet_id),
					Bookinglist_mtick_id:      ticket.Mtick_id,
					Bookinglist_mid:           grup.Group_mid,
				}

				db.DB[0].NewRecord(bookdet)

				if err := db.DB[0].Create(&booklist).Error; err != nil {
					db.DB[0].Where("bookinglist_bookingdet_id = ?", bookdet.Bookingdet_id).Delete(entities.Bookinglist{})
					db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Delete(entities.Bookingdet{})
					db.DB[0].Where("booking_id = ?", booking.Booking_id).Delete(entities.Booking{})

					return nil, "07", "Error when inserting bookingdet data (" + err.Error() + ")", false
				}
			}
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.Email)
	m.SetHeader("To", r.Email)
	m.SetHeader("Subject", "Booking detail")
	m.SetBody("text/html", `
	<html>
  		<head>
		<style>
			.coupon {
			border: 5px dotted #bbb; /* Dotted border */
			width: 80%;
			border-radius: 15px; /* Rounded border */
			margin: 0 auto; /* Center the coupon */
			max-width: 600px;
			}

			.container {
			padding: 2px 16px;
			background-color: #f1f1f1;
			}

			.promo {
			background: #ccc;
			padding: 3px;
			}

			.expire {
			color: red;
			}
		</style>
		</head>
		<body>
		<div class="coupon">
		<div class="container">
			<center><h3>AINO Indonesia</h3></center>
		</div>
		<div class="container" style="background-color:white">
			<p>Congratulation! Your booking detail are listed below:</p>
		</div>
		<div class="container">
			<p>UUID: `+booking.Booking_uuid+`</p>
			<p>Booking number: <span class="promo">`+booking.Booking_number+`</span></p>
		</div>
		</div>
		</body>
	</html>
	`)

	dialer := gomail.NewPlainDialer(
		config.Mail.Host,
		config.Mail.Port,
		config.Mail.Username,
		config.Mail.Password,
	)

	err := dialer.DialAndSend(m)
	if err != nil {
		return map[string]interface{}{
			"data_tariff": r.Trf,
			"booking_detail": map[string]interface{}{
				"booking_UUID":   booking.Booking_uuid,
				"booking_number": booking.Booking_number,
			},
		}, "08", "Booking success, but an error occurred when sending e-mail (" + err.Error() + ")", true
	}

	return map[string]interface{}{
		"data_tariff": r.Trf,
		"booking_detail": map[string]interface{}{
			"booking_UUID":   booking.Booking_uuid,
			"booking_number": booking.Booking_number,
		},
	}, "01", "Booking success, email sent", true
}
