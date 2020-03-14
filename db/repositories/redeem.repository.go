package repositories

import (
	"strings"
	"twc-ota-api/config"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"

	"gopkg.in/gomail.v2"
)

// RedeemTicket : redeem ticket
func RedeemTicket(token *entities.Users, r *requests.RedeemReq) (map[string]interface{}, string, string, bool) {

	qrCode := "TIC" + "921a924e-4df0-4296-bab3-59a533d33eb4#1559722940"
	qrString := strings.ReplaceAll(qrCode, "#", "%23")
	qrImage := `<img src="https://chart.apis.google.com/chart?cht=qr&chs=300x300&chl=` + qrString + `&chld=H|0" />`

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
		return nil, "08", "Redeem success, but an error occurred when sending e-mail (" + err.Error() + ")", true
	}
	return map[string]interface{}{
		"request": r,
	}, "01", "Redeem success, email sent", true
}
