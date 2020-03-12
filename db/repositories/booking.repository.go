package repositories

import (
	"twc-ota-api/config"
	"twc-ota-api/db/entities"

	"gopkg.in/gomail.v2"
)

// const CONFIG_SMTP_HOST = "smtp.gmail.com"
// const CONFIG_SMTP_PORT = 587
// const CONFIG_EMAIL = "zurinezekiel@gmail.com"
// const CONFIG_USERNAME = "zurinezekiel@gmail.com"
// const CONFIG_PASSWORD = "immortalblood"

// BookingTicket : booking ticket
func BookingTicket(token *entities.Users) (map[string]interface{}, string, string, bool) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.Email)
	m.SetHeader("To", token.Email)
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
			<p>UUID: jhaw-78hajd-878gjw-823jn</p>
			<p>Booking number: <span class="promo">18237034</span></p>
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
		return nil, "02", "Error when sending e-mail (" + err.Error() + ")", false
	}

	return map[string]interface{}{
		"data": token,
	}, "01", "Booking success, email sent", true
}
