package repositories

import (
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/middleware"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//GetUser : select data user
func GetUser(r interface{}) (map[string]interface{}, string, string, bool) {
	var user entities.Users

	email := r.(map[string]interface{})["email"]
	password := r.(map[string]interface{})["password"]

	if email == nil || password == nil {
		return nil, "99", "Email or password cant't be empty", false
	}

	if err := db.DB[0].Where("email = ? AND (type = 'AT' OR type = 'TRPPLNR' or type = 'B2B')", email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", "Email not registered (" + err.Error() + ")", false
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

	db.DB[1].Where("email = ?", r.Email).Find(&user)

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
