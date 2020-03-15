package repositories

import (
	"strconv"
	"strings"
	"time"
	"twc-ota-api/config"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gomail.v2"
)

// RedeemTicket : redeem ticket
func RedeemTicket(token *entities.Users, r *requests.RedeemReq) (map[string]interface{}, string, string, bool) {
	if r.BookNumber == "" {
		return nil, "99", "Format error, booking number is required", false
	}

	var qrImage string
	var dataTrf []entities.TrfResponse

	var bookings []entities.Booking
	if err := db.DB[0].Where("booking_number = ?", r.BookNumber).Find(&bookings).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Booking data not found (" + err.Error() + ")", false
	}

	if len(bookings) == 0 {
		return nil, "02", "Booking data not found", false
	}

	for _, booking := range bookings {
		stan := time.Now().UnixNano()
		tickID := uuid.NewV4()
		tick := entities.TickModel{
			Tick_id:             tickID,
			Tick_amount:         booking.Booking_amount,
			Tick_date:           booking.Booking_date,
			Tick_emoney:         booking.Booking_emoney,
			Tick_issuing:        booking.Booking_date,
			Tick_mid:            booking.Booking_mid,
			Tick_number:         booking.Booking_number,
			Tick_payment_method: booking.Booking_payment_method,
			Tick_purc:           time.Now().Format("2006-01-02 15:04:05"),
			Tick_src_type:       5,
			Tick_total_payment:  booking.Booking_total_payment,
			Tick_stan:           int8(stan),
		}
		db.DB[0].NewRecord(tick)

		if err := db.DB[0].Create(&tick).Error; err != nil {
			return nil, "03", "Error when inserting ticket data (" + err.Error() + ")", false
		}

		var bookingdets []entities.Bookingdet

		if err := db.DB[0].Where("bookingdet_booking_id = ?", booking.Booking_id).Find(&bookingdets).Error; gorm.IsRecordNotFoundError(err) {
			db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickDetModel{})
			db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
			return nil, "04", "Bookingdet data not found (" + err.Error() + ")", false
		}

		if len(bookingdets) == 0 {
			db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
			db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
			return nil, "04", "Bookingdet data not found", false
		}

		for _, bookingdet := range bookingdets {
			tickdetID := uuid.NewV4()
			qrCode := tickdetID.String() + "#" + strconv.FormatInt(stan, 10)

			tickdet := entities.TickDetModel{
				Tickdet_id:      tickdetID,
				Tickdet_tick_id: tick.Tick_id,
				Tickdet_amount:  bookingdet.Bookingdet_amount,
				Tickdet_qr:      qrCode,
				Tickdet_qty:     bookingdet.Bookingdet_qty,
				Tickdet_total:   bookingdet.Bookingdet_total,
				Tickdet_trf_id:  bookingdet.Bookingdet_trf_id,
				Tickdet_trftype: bookingdet.Bookingdet_trftype,
			}
			db.DB[0].NewRecord(tickdet)

			if err := db.DB[0].Create(&tickdet).Error; err != nil {
				db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickDetModel{})
				db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
				return nil, "05", "Error when inserting ticketdet data (" + err.Error() + ")", false
			}

			var bookinglists []entities.Bookinglist
			if err := db.DB[0].Where("bookinglist_bookingdet_id = ?", bookingdet.Bookingdet_id).Find(&bookinglists).Error; gorm.IsRecordNotFoundError(err) {
				db.DB[0].Where("ticklist_tickdet_id = ?", tickdetID).Delete(entities.TickListModel{})
				db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickDetModel{})
				db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
				return nil, "06", "Bookinglist data not found (" + err.Error() + ")", false
			}

			if len(bookinglists) == 0 {
				db.DB[0].Where("ticklist_tickdet_id = ?", tickdetID).Delete(entities.TickListModel{})
				db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickDetModel{})
				db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
				return nil, "06", "Bookinglist data not found", false
			}

			for _, bookinglist := range bookinglists {
				ticklistID := uuid.NewV4()

				ticklist := entities.TickListModel{
					Ticklist_id:         ticklistID,
					Ticklist_mid:        bookinglist.Bookinglist_mid,
					Ticklist_tickdet_id: tickdetID,
					Ticklist_mtick_id:   bookinglist.Bookinglist_mtick_id,
				}

				db.DB[0].NewRecord(ticklist)

				if err := db.DB[0].Create(&ticklist).Error; err != nil {
					db.DB[0].Where("ticklist_tickdet_id = ?", tickdetID).Delete(entities.TickListModel{})
					db.DB[0].Where("tickdet_tick_id = ?", tick.Tick_id).Delete(entities.TickDetModel{})
					db.DB[0].Where("tick_id = ?", tick.Tick_id).Delete(entities.TickModel{})
					return nil, "07", "Error when inserting ticketlist data (" + err.Error() + ")", false
				}
			}

			tmpTrf := entities.TrfResponse{
				TickDetID:      tickdetID,
				TickDetAmount:  tickdet.Tickdet_amount,
				TickDetQty:     tickdet.Tickdet_qty,
				TickDetTrfID:   tickdet.Tickdet_trf_id,
				TickDetTrfType: tickdet.Tickdet_trftype,
				TickDetQr:      "TIC" + qrCode,
			}

			dataTrf = append(dataTrf, tmpTrf)
			qrString := strings.ReplaceAll(qrCode, "#", "%23")
			qrImage = qrImage + `<img src="https://chart.apis.google.com/chart?cht=qr&chs=300x300&chl=` + `TIC` + qrString + `&chld=H|0" />`
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.Email)
	m.SetHeader("To", "rinoridlojulianto@gmail.com")
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
			<p>Congratulation on your successful redeem. Here is your QRCode:</p>
		</div>
		<div class="container">
			<center>`+qrImage+`</center>
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
			"trf_data": dataTrf,
		}, "08", "Redeem success, but an error occurred when sending e-mail (" + err.Error() + ")", true
	}
	return map[string]interface{}{
		"trf_data": dataTrf,
	}, "01", "Redeem success, email sent", true
}
