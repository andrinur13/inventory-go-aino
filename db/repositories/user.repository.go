package repositories

import (
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

	if err := db.DB[0].Where("email = ? AND (type = 'AT' OR type = 'TRPPLNR')", email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
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
