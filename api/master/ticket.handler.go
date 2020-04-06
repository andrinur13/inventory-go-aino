package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

var cm *service.CacheManager

// TicketRouter : Routing
func TicketRouter(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	ticket := r.Group("/ticket")
	{
		ticket.POST("/list", permission.Set("PERMISSION_MASTER_USER_VIEW", GetTicketList))
		ticket.POST("/booking", permission.Set("PERMISSION_MASTER_USER_SAVE", BookingTicket))
		ticket.POST("/redeem", permission.Set("PERMISSION_MASTER_USER_SAVE", RedeemTicket))
		ticket.POST("/checkout", permission.Set("PERMISSION_MASTER_USER_SAVE", CheckoutTicket))
		ticket.POST("/tes", permission.Set("PERMISSION_MASTER_USER_VIEW", Tes))
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

//Tes : for testing purpose
func Tes(c *gin.Context) {
	nano := time.Now().UnixNano()
	unix := time.Now().Unix()
	micro := nano / (int64(time.Millisecond) / int64(time.Nanosecond))

	data := map[string]interface{}{
		"second": unix,
		"micro":  micro,
		"nano":   nano,
		"dummy":  "This is dummy data",
	}
	c.JSON(http.StatusOK, builder.ApiResponse(true, "Success testing auth", "01", data))
}
