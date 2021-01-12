package repositories

import (
	"encoding/json"
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"
)

// InsertFav : insert to table favorite
func InsertFav(token *entities.Users, r *requests.FavReq) (map[string]interface{}, string, string, bool) {
	if len(r.Data) == 0 {
		return nil, "99", "Favorite data is required", false
	}

	// if r.TotalAmount == 0 {
	// 	return nil, "99", "Total amount is required", false
	// }
	if r.Bruto == 0 {
		return nil, "99", "Bruto is required", false
	}

	if r.Netto == 0 {
		return nil, "99", "Netto is required", false
	}

	if r.Name == "" {
		return nil, "99", "Favorite's name is required", false
	}

	rData, err := json.Marshal(&r.Data)
	if err != nil {
		return nil, "99", "Failed to parse json key data (" + err.Error() + ")", false
	}

	jData := string(rData)

	extras := `{"name":"` + r.Name + `",` + `"price_bruto":` + strconv.FormatFloat(r.Bruto, 'f', -1, 32) + `, "price_disc":` + strconv.FormatFloat(r.Disc, 'f', -1, 32) + `, "price_netto":` + strconv.FormatFloat(r.Netto, 'f', -1, 32) + `, "data":` + jData + `}`

	fav := entities.Favorite{
		Fav_user_id: token.ID,
		Fav_data:    extras,
		Fav_created: time.Now().Format("2006-01-02 15:04:05"),
	}

	db.DB[0].NewRecord(fav)

	if err := db.DB[0].Create(&fav).Error; err != nil {
		return nil, "01", "Error when inserting favorite data (" + err.Error() + ")", false
	}

	return nil, "00", "Success insert favorite", true
}
