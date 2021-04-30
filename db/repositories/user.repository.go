package repositories

import (
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"fmt"
	"reflect"
	"time"
	"twc-ota-api/config"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/utils/helper"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//GetUser : select data user
func GetUser(r interface{}) (map[string]interface{}, string, string, bool) {
	var user entities.Users

	email := r.(map[string]interface{})["email"]
	password := r.(map[string]interface{})["password"]

	if email == nil || password == nil {
		return nil, "99", "Email or password cant't be empty", false
	}

	erro := db.DB[0].Where("email = ? AND (type = 'AT' OR type = 'TRPPLNR' or type = 'B2B')", email).Find(&user).Error;

	//If Connection Refused
	if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
		fmt.Printf("%v \n", erro.Error())
			for i := 0; i<4; i++ {
				erro = db.DB[0].Where("email = ? AND (type = 'AT' OR type = 'TRPPLNR' or type = 'B2B')", email).Find(&user).Error;
				if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
					fmt.Printf("Hitback(%d)%v \n", i, erro)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	if gorm.IsRecordNotFoundError(erro) {
		return nil, "02", "Email not registered (" + erro.Error() + ")", false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.(string)))

	if err != nil {
		return nil, "03", "Invalid credential", false
	}

	token, _ := middleware.CreateJwtToken(&user)

	return map[string]interface{}{
		"email":    email,
		"agent_id": user.Typeid,
		"token":    token,
	}, "01", "ok", true
}

//InsertUser : insert data user
func InsertUser(r entities.UserReq) (map[string]interface{}, string, string, bool) {
	if r.Name == "" {
		return nil, "99", "Name cant't be empty", false
	}

	if r.Email == "" {
		return nil, "99", "E-mail cant't be empty", false
	}

	if r.Password == "" {
		return nil, "99", "Password cant't be empty", false
	}

	if r.NationalityID == 0 {
		return nil, "99", "Nationality ID cant't be empty", false
	}

	if r.Type == "" {
		return nil, "99", "User's type cant't be empty", false
	}

	var user []entities.Users

	erro := db.DB[1].Where("email = ?", r.Email).Find(&user).Error;

	//If Connection refused
	if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
		fmt.Printf("%v \n", erro.Error())
		fmt.Printf("%v \n", reflect.TypeOf(erro).String())
			for i := 0; i<4; i++ {
				erro = db.DB[1].Where("email = ?", r.Email).Find(&user).Error;
				if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError") {
					fmt.Printf("Hitback(%d)%v \n", i, erro)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	// db.DB[1].Where("email = ?", r.Email).Find(&user)

	if len(user) > 0 {
		return nil, "02", "E-mail already registered", false
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "99", "General error, failed hashing password", false
	}

	var dataUser entities.Users

	if r.Typeid != 0 {
		dataUser = entities.Users{
			Name:         r.Name,
			Email:        r.Email,
			Password:     string(password),
			Type:         r.Type,
			Users_extras: `{"verified": true, "nationality_id": ` + strconv.Itoa(r.NationalityID) + `}`,
			Typeid:       r.Typeid,
			Created_at:   time.Now().Format("2006-01-02 15:04:05"),
		}
	} else {
		dataUser = entities.Users{
			Name:         r.Name,
			Email:        r.Email,
			Password:     string(password),
			Type:         r.Type,
			Users_extras: `{"verified": true, "nationality_id": ` + strconv.Itoa(r.NationalityID) + `}`,
			Created_at:   time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	db.DB[1].NewRecord(dataUser)

	if err := db.DB[1].Create(&dataUser).Error; err != nil {
		return nil, "03", "Error when inserting user data (" + err.Error() + ")", false
	}

	return map[string]interface{}{
		"name_depan": r.Name,
		"email":      r.Email,
		"type":       r.Type,
	}, "01", "Registration success", true
}

//UpdatePassword : update password user
func UpdatePassword(token *entities.Users, r *requests.UpdatePassword) (map[string]interface{}, string, string, bool) {
	var user entities.Users

	if r.OldPwd == "" {
		return nil, "99", "Old Password is required", false
	}

	if r.NewPwd == "" {
		return nil, "99", "New Password is required", false
	}

	if r.ConfPwd == "" {
		return nil, "99", "Confirm Password is required", false
	}

	if r.NewPwd != r.ConfPwd {
		return nil, "02", "Confirmation password doesn't match", false
	}

	erro := db.DB[0].Where("id = ? AND (type = 'AT' OR type = 'TRPPLNR' or type = 'B2B')", token.ID).Find(&user).Error;

	//If Connection refused
	if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
		fmt.Printf("%v \n", erro.Error())
		fmt.Printf("%v \n", reflect.TypeOf(erro).String())
			for i := 0; i<4; i++ {
				erro = db.DB[0].Where("id = ? AND (type = 'AT' OR type = 'TRPPLNR' or type = 'B2B')", token.ID).Find(&user).Error;
				if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError") {
					fmt.Printf("Hitback(%d)%v \n", i, erro)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	if gorm.IsRecordNotFoundError(erro) {
		return nil, "03", "User not found (" + erro.Error() + ")", false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPwd))

	if err != nil {
		return nil, "04", "Wrong old password", false
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.NewPwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, "05", "General error, failed hashing password", false
	}

	if err := db.DB[0].Model(&user).Where("id = ?", token.ID).Update("password", string(password)).Error; err != nil {
		return nil, "06", "Failed to update password", false
	}

	return map[string]interface{}{
		"user_id":  token.ID,
		"email":    token.Email,
		"agent_id": user.Typeid,
	}, "01", "Password succesfully updated", true
}

//ResetPassword : reset forgoten password user
func ResetPassword(r *requests.ResetPassword) (map[string]interface{}, string, string, bool) {

	if r.Email == "" {
		return nil, "99", "Email is required", false
	}

	if e := r.Email; !helper.IsEmailValid(e) {
		return nil, "99", "Invalid email address!", false
	}

	var user entities.Users
	var passRes entities.PasswordReset

	erro := db.DB[1].Where("email = ?", r.Email).Find(&user).Error;

	//If Connection refused
	if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
		fmt.Printf("%v \n", erro.Error())
		fmt.Printf("%v \n", reflect.TypeOf(erro).String())
			for i := 0; i<4; i++ {
				erro = db.DB[1].Where("email = ?", r.Email).Find(&user).Error;
				if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError") {
					fmt.Printf("Hitback(%d)%v \n", i, erro)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (erro != nil) && (reflect.TypeOf(erro).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	if gorm.IsRecordNotFoundError(erro) {
		return nil, "02", "E-mail not registered! (" + erro.Error() + ")", false
	}

	db.DB[1].Where("email = ?", r.Email).Delete(&passRes)

	char := helper.StringWithCharset(14, charset)

	hasher := sha1.New()
	hasher.Write([]byte(char))
	token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	var pwRst entities.PasswordReset

	pwRst = entities.PasswordReset{
		Email:      r.Email,
		Token:      string(token),
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}

	db.DB[1].NewRecord(pwRst)

	if err := db.DB[1].Create(&pwRst).Error; err != nil {
		return nil, "03", "Error when inserting data to table reset password (" + err.Error() + ")", false
	}

	url := config.App.GtHost + "password/reset?email=" + r.Email + "&token=" + string(token)

	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.Email)
	m.SetHeader("To", r.Email)
	m.SetHeader("Subject", "AINO Indonesia | Reset Password")
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
			<p>You were requested for user's password reset on our platform. Please see the detail below and access the given url:</p>
		</div>
		<div class="container">
			<p>Email: `+r.Email+`</p>
			<p>Reset password URL: <span class="promo"><a href="`+url+`">`+url+`</a></span></p>
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
			"email": r.Email,
		}, "03", "Error occured, failed to sent mail (" + err.Error() + ")", false
	}
	return map[string]interface{}{
		"email": r.Email,
	}, "01", "Password reset url has been sent to your e-mail", true
}

//UpdateResetPassword : update reset password user
func UpdateResetPassword(r *requests.UpdateResetPassword) (map[string]interface{}, string, string, bool) {
	if r.Email == "" {
		return nil, "99", "E-mail is required", false
	}

	if e := r.Email; !helper.IsEmailValid(e) {
		return nil, "99", "Invalid email address!", false
	}

	if r.Token == "" {
		return nil, "99", "Token is required", false
	}

	if r.NewPwd == "" {
		return nil, "99", "Password is required", false
	}

	if r.ConfPwd == "" {
		return nil, "99", "Confirm Password is required", false
	}

	if r.NewPwd != r.ConfPwd {
		return nil, "02", "Confirmation password doesn't match", false
	}

	var user entities.Users
	var rst entities.PasswordReset

	err := db.DB[1].Where("email = ?", r.Email).Find(&user).Error;

	//If Connection refused
	if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
		fmt.Printf("%v \n", err.Error())
		fmt.Printf("%v \n", reflect.TypeOf(err).String())
			for i := 0; i<4; i++ {
				err = db.DB[1].Where("email = ?", r.Email).Find(&user).Error;
				if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError") {
					fmt.Printf("Hitback(%d)%v \n", i, err)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
		if (err != nil) && (reflect.TypeOf(err).String() == "*net.OpError"){
			return nil, "502", "Connection has a problem", false
		}
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, "03", "E-mail not registered! (" + err.Error() + ")", false
	}

	if err := db.DB[1].Where("email = ? AND token = ?", r.Email, r.Token).Find(&rst).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "04", "Password reset request not found! Please repeat your password reset step. (" + err.Error() + ")", false
	}

	if err := db.DB[1].Where("email = ?", r.Email).Delete(&rst).Error; err != nil {
		return nil, "05", "Failed to delete reset data. (" + err.Error() + ")", false
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.NewPwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, "06", "General error, failed hashing password", false
	}

	if err := db.DB[0].Model(&user).Where("email = ?", r.Email).Update("password", string(password)).Error; err != nil {
		return nil, "07", "Failed to update password", false
	}

	return map[string]interface{}{
		"user_id":  user.ID,
		"email":    r.Email,
		"agent_id": user.Typeid,
	}, "01", "Password succesfully updated", true
}
