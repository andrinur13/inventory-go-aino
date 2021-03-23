package master

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"twc-ota-api/config"
	"twc-ota-api/db/entities"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gabriel-vasile/mimetype"
	// limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const IMGPREFIX string = "FAVIMAGE:"
const MIMETYPEJPEG string = "image/jpeg"
const MIMETYPEPNG string = "image/png"

var cm *service.CacheManager

// TicketRouter : Routing
func TicketRouter(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	// r.Use(limits.RequestSizeLimiter(2000))
	cm = cacheManager
	ticket := r.Group("/ticket")
	{
		ticket.POST("/list", permission.Set("PERMISSION_MASTER_USER_VIEW", GetTicketList))
		ticket.POST("/booking", permission.Set("PERMISSION_MASTER_USER_SAVE", BookingTicket))
		ticket.POST("/redeem", permission.Set("PERMISSION_MASTER_USER_SAVE", RedeemTicket))
		ticket.POST("/checkout", permission.Set("PERMISSION_MASTER_USER_SAVE", CheckoutTicket))
		ticket.GET("/cluster", permission.Set("PERMISSION_MASTER_USER_VIEW", GetCluster))
		ticket.POST("/tes", permission.Set("PERMISSION_MASTER_USER_VIEW", Tes))
	}
	site := r.Group("/site")
	{
		site.GET("/detail", permission.Set("PERMISSION_MASTER_USER_VIEW", DetailSite))
	}
	agent := r.Group("/register")
	{
		agent.POST("/agent", permission.Set("PERMISSION_MASTER_USER_SAVE", RegisterAgent))
	}
	trx := r.Group("/trx")
	{
		trx.GET("list/:page/:size", permission.Set("PERMISSION_MASTER_USER_VIEW", GetTransaction))
		trx.POST("/create", permission.Set("PERMISSION_MASTER_USER_SAVE", CreateTrx))
		trx.PUT("/update", permission.Set("PERMISSION_MASTER_USER_SAVE", UpdateTrx))
		trx.PUT("/pay", permission.Set("PERMISSION_MASTER_USER_SAVE", UpdateTrxPay))
		trx.POST("/info", permission.Set("PERMISSION_MASTER_USER_SAVE", GetInfo))
		trx.POST("/number", permission.Set("PERMISSION_MASTER_USER_SAVE", GetNumber))
	}
	fav := r.Group("/fav")
	{
		fav.POST("/create", permission.Set("PERMISSION_MASTER_USER_SAVE", CreateFav))
		fav.POST("/delete", permission.Set("PERMISSION_MASTER_USER_SAVE", DeleteFav))
		fav.GET("/list", permission.Set("PERMISSION_MASTER_USER_SAVE", GetFav))
		fav.POST("/image", permission.Set("PERMISSION_MASTER_USER_SAVE", UploadFavImage))
	}
}

// GetTicketList : Get ticket's fare data
func GetTicketList(c *gin.Context) {
	var param interface{}
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetTicket(param, userData)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	//test timeout client
	// time.Sleep(3 * time.Second)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// RedeemTicket : redeem ticket
func RedeemTicket(c *gin.Context) {
	var param requests.RedeemReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.RedeemTicket(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// BookingTicket : booking ticket
func BookingTicket(c *gin.Context) {
	var param requests.BookingReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.BookingTicket(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// CheckoutTicket : checkout ticket
func CheckoutTicket(c *gin.Context) {
	var param entities.CheckOutReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.CheckoutB2B(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// GetCluster : get data cluster
func GetCluster(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	nationality := c.DefaultQuery("nationality", "")

	data, code, msg, stat := repositories.SelectCluster(userData, nationality)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

// DetailSite : get data site detail
func DetailSite(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	nationality := c.DefaultQuery("nationality_id", "")
	siteID := c.DefaultQuery("site_id", "")

	data, code, msg, stat := repositories.GetSite(userData, nationality, siteID)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

// RegisterAgent : register agent
func RegisterAgent(c *gin.Context) {
	var param entities.AgentReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.InsertAgent(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// GetTransaction : Get transaction data
func GetTransaction(c *gin.Context) {
	status := strings.ToUpper(c.Query("status"))
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusOK, builder.ApiResponse(false, "General Error, couldn't parse page", "99", nil))
		logger.Warning("General Error, couldn't parse page", "99", false, "")
		return
	}
	size, err := strconv.Atoi(c.Param("size"))
	if err != nil {
		c.JSON(http.StatusOK, builder.ApiResponse(false, "General Error, couldn't parse size", "99", nil))
		logger.Warning("General Error, couldn't parse size", "99", false, "")
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat, totalData, totalPages, currentData := repositories.SelectTrip(userData, page, size, status)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	//test timeout client
	// time.Sleep(3 * time.Second)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	// c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	c.JSON(http.StatusOK, builder.ListResponse(stat, msg, code, totalData, currentData, totalPages, page, size, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

// CreateTrx : create new trx
func CreateTrx(c *gin.Context) {
	var param requests.TrxReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.InsertTrx(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// UpdateTrx : update trx
func UpdateTrx(c *gin.Context) {
	var param requests.TrxReqUpdate
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.UpdateTrx(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// UpdateTrxPay : update trx payment
func UpdateTrxPay(c *gin.Context) {
	var param requests.TrxReqUpdate
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.UpdateTrxPayment(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// GetInfo : get ticket info from invoice
func GetInfo(c *gin.Context) {
	var param requests.TrxQReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetQR(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// GetNumber : get trx by number
func GetNumber(c *gin.Context) {
	var param requests.TrxQReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetTrxByNumber(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// CreateFav : create favorite
func CreateFav(c *gin.Context) {
	var param requests.FavReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	var filePath string
	var imageIsStored bool
	var stat bool
	data := make(map[string]interface{})
	code := "01"
	msg := "Failed to save image"
	//generate UUID
	favID := uuid.NewV4()
	if param.Image != "" {
		//convert base64 to bytearray
		if byteImage, e := base64.StdEncoding.DecodeString(param.Image); e != nil {
			imageIsStored = false
		} else {
			//check for uploaded file mimetype
			allowedMIMETypes := []string{MIMETYPEJPEG, MIMETYPEPNG}
			mime := mimetype.Detect(byteImage)
			if !mimetype.EqualsAny(mime.String(), allowedMIMETypes...) {
				imageIsStored = false
			} else {
				//generate file name
				hashImageName := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v%v", IMGPREFIX, favID))))
				filePath = filepath.Join(config.App.ImageDirectory, filepath.Base(hashImageName))
				//add file extension
				switch mime.String() {
				case MIMETYPEJPEG:
					{
						filePath += ".jpg"
					}
				case MIMETYPEPNG:
					{
						filePath += ".png"
					}
				default:
					{
						//don't add anything
					}
				}
				//try removing file
				os.Remove(filePath)
				//write image file to storage
				if file, e := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666); e != nil {
					imageIsStored = false
				} else {
					defer file.Close()
					if _, e := file.Write(byteImage); e != nil {
						imageIsStored = false
					} else {
						if e := file.Sync(); e != nil {
							imageIsStored = false
						} else {
							imageIsStored = true
							rgx := regexp.MustCompile("/var/www/html/public")
							param.ImageURL = rgx.ReplaceAllString(filePath, "/static")
						}
					}
				}
			}
		}
	} else {
		imageIsStored = true
	}
	if imageIsStored {
		data, code, msg, stat = repositories.InsertFav(favID, userData, &param)
	}

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// DeleteFav : create favorite
func DeleteFav(c *gin.Context) {
	var param requests.FavDelete
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.DeleteFav(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// GetFav : get favorite data
func GetFav(c *gin.Context) {
	in, _ := json.Marshal(nil)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.SelectFav(userData)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

//UploadFavImage :
func UploadFavImage(c *gin.Context) {
	var param requests.FavUploadImage
	in, _ := json.Marshal(param)
	var success bool
	var message string
	code := "01"
	if e := c.BindJSON(&param); e != nil {
		message = "Invalid request format"
	} else {
		//convert base64 to bytearray
		if byteImage, e := base64.StdEncoding.DecodeString(param.Image); e != nil {
			message = "Unable to decode base64 image"
		} else {
			//check for uploaded file mimetype
			allowedMIMETypes := []string{MIMETYPEJPEG, MIMETYPEPNG}
			mime := mimetype.Detect(byteImage)
			if !mimetype.EqualsAny(mime.String(), allowedMIMETypes...) {
				message = "Image format not supported, please use either jpg or png image format."
			} else {
				//generate file name
				hashImageName := fmt.Sprintf("%x", sha256.Sum256([]byte(IMGPREFIX+param.FavID)))
				filePath := filepath.Join(config.App.ImageDirectory, filepath.Base(hashImageName))
				//add file extension
				switch mime.String() {
				case MIMETYPEJPEG:
					{
						filePath += ".jpg"
					}
				case MIMETYPEPNG:
					{
						filePath += ".png"
					}
				default:
					{
						//don't add anything
					}
				}
				//try removing file
				os.Remove(filePath)
				//write image file to storage
				if file, e := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666); e != nil {
					message = fmt.Sprintf("Unable to store image (%v)", e.Error())
				} else {
					defer file.Close()
					if _, e := file.Write(byteImage); e != nil {
						message = fmt.Sprintf("Unable to store image (%v)", e.Error())
					} else {
						if e := file.Sync(); e != nil {
							message = fmt.Sprintf("Unable to store image (%v)", e.Error())
						} else {
							//update to database
							rgx := regexp.MustCompile("/var/www/html/public")
							if e := repositories.StoreFavImage(param.FavID, rgx.ReplaceAllString(filePath, "/static")); e != nil {
								//remove file if update file
								os.Remove(filePath)
								message = fmt.Sprintf("Unable to store image (%v)", e.Error())
							} else {
								code = "00"
								message = "Image stored successfully"
							}
						}
					}
				}
			}
		}
	}

	out, _ := json.Marshal(gin.H{})
	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(success, message, code, gin.H{}))
	logger.Info(message, code, success, fmt.Sprintf("%v", gin.H{}), string(in))
}

//Tes : for testing purpose
func Tes(c *gin.Context) {
	nano := time.Now().UnixNano()
	unix := time.Now().Unix()
	micro := nano / (int64(time.Millisecond) / int64(time.Nanosecond))

	t, _ := time.Parse("2006-01-02", "2020-06-30")
	dayExp := t.Add(time.Hour*time.Duration((2*24))).Format("2006-01-02") + " 23:59:59"

	data := map[string]interface{}{
		"day_exp": dayExp,
		"second":  unix,
		"micro":   micro,
		"nano":    nano,
		"dummy":   "This is dummy data",
	}
	c.JSON(http.StatusOK, builder.ApiResponse(true, "Success testing auth", "01", data))
}
