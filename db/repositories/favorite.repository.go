package repositories

import (
	"encoding/json"
	"strconv"
	"time"
	"twc-ota-api/db"
	"twc-ota-api/db/entities"
	"twc-ota-api/requests"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// InsertFav : insert to table favorite
func InsertFav(favID uuid.UUID, token *entities.Users, r *requests.FavReq) (map[string]interface{}, string, string, bool) {
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

	extras := `{"name":"` + r.Name + `", "image_url":"` + r.ImageURL + `", "duration":` + strconv.Itoa(r.Duration) + `, "adult":` + strconv.Itoa(r.Adult) + `, "child":` + strconv.Itoa(r.Child) + `, "nationality_id":` + strconv.Itoa(r.NationalityID) + `, "price_bruto":` + strconv.FormatFloat(r.Bruto, 'f', -1, 32) + `, "price_disc":` + strconv.FormatFloat(r.Disc, 'f', -1, 32) + `, "price_netto":` + strconv.FormatFloat(r.Netto, 'f', -1, 32) + `, "data":` + jData + `}`

	fav := entities.Favorite{
		Fav_id:      favID,
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

// SelectFav : select from table favorite
func SelectFav(token *entities.Users) (*[]entities.FavResp, string, string, bool) {
	var fav []entities.Favorite

	if err := db.DB[0].Select("fav_id, fav_data").Where("fav_deleted is null and fav_user_id = ?", token.ID).Order("fav_created desc").Find(&fav).Error; err != nil {
		return nil, "01", err.Error(), false
	}

	if len(fav) == 0 {
		return nil, "02", "Favorite data not found", false
	}

	var resp []entities.FavResp

	for _, data := range fav {
		var jParse requests.FavReq
		json.Unmarshal([]byte(data.Fav_data), &jParse)

		var favData []entities.FavData

		for _, dataFav := range jParse.Data {

			var favTrf []entities.FavTrf
			for _, trfFav := range dataFav.Trf {
				var trf entities.TariffModel

				if err := db.DB[1].Select(`trf_name`).Where("deleted_at is null and trf_id = ?", trfFav.TrfID).Find(&trf).Error; gorm.IsRecordNotFoundError(err) {
					return nil, "03", "No tarif data found (" + err.Error() + ")", false
				}

				tmpTrfFav := entities.FavTrf{
					TrfID:    trfFav.TrfID,
					TrfCode:  trfFav.TrfCode,
					TrfQty:   trfFav.TrfQty,
					TrfNetto: trfFav.TrfNetto,
					TrfName:  trf.Trf_name,
				}

				favTrf = append(favTrf, tmpTrfFav)
			}

			tmpFavData := entities.FavData{
				Day: dataFav.Day,
				Trf: favTrf,
			}

			favData = append(favData, tmpFavData)
		}

		var image_url string

		image_url = jParse.ImageURL

		if image_url == "" {
			image_url = "static/b2bm/package/default_package.jpg"
		}

		tmpResp := entities.FavResp{
			Name:          jParse.Name,
			ImageURL:      image_url,
			Duration:      jParse.Duration,
			NationalityID: jParse.NationalityID,
			Adult:         jParse.Adult,
			Child:         jParse.Child,
			Bruto:         jParse.Bruto,
			Netto:         jParse.Netto,
			Disc:          jParse.Disc,
			Data:          favData,
			PaketID:       data.Fav_id,
		}

		resp = append(resp, tmpResp)
	}

	return &resp, "00", "Success get favorite", true
}

// DeleteFav : delete from table favorite
func DeleteFav(token *entities.Users, r *requests.FavDelete) (map[string]interface{}, string, string, bool) {
	if r.FavID == "" {
		return nil, "99", "ID is required", false
	}

	var fav entities.Favorite

	if err := db.DB[0].Select(`fav_user_id`).Where("fav_deleted is null and fav_id = ?", r.FavID).Find(&fav).Error; gorm.IsRecordNotFoundError(err) {
		return nil, "02", err.Error(), false
	}

	if fav.Fav_user_id != token.ID {
		return nil, "03", "Couldn't delete favorite, invalid privileges", false
	}

	if err := db.DB[0].Model(&fav).Where("fav_id = ?", r.FavID).Update("fav_deleted", time.Now().Format("2006-01-02 15:04:05")).Error; err != nil {
		return nil, "01", "Error when deleting favorite data (" + err.Error() + ")", false
	}

	return nil, "00", "Success delete favorite", true
}

//StoreFavImage :
func StoreFavImage(favID, imagePath string) error {
	if e := db.DB[0].Exec(`
		UPDATE public.favorite
		SET fav_data = jsonb_set(fav_data, '{image_url}', '"`+imagePath+`"')
		WHERE fav_id = ?;
	`, favID).GetErrors(); len(e) > 0 {
		return e[0]
	}
	return nil
}
