package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

var cm *service.CacheManager

// Param : parameter for ticket list
type Param struct {
	Mbmid string `json:"merchant_code" binding:"required"`
	Mbtid string `json:"device_code"`
	Ctgid int    `json:"ctg_id"`
}

// TicketRouter : Routing
func TicketRouter(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	ticket := r.Group("/ticket")
	{
		ticket.POST("/list", permission.Set("PERMISSION_MASTER_USER_VIEW", GetTicketList))
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

//Tes : for testing purpose
func Tes(c *gin.Context) {
	data := map[string]interface{}{
		"dummy": "This is dummy data",
	}
	c.JSON(http.StatusOK, builder.ApiResponse(true, "Success testing auth", "01", data))
}
