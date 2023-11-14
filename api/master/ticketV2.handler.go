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
	"twc-ota-api/requests"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

// TicketRouter : Routing
func TicketRouterV2(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	ticket := r.Group("/ticket")
	{
		ticket.POST("/redeem", permission.Set("PERMISSION_MASTER_USER_SAVE", RedeemTicketV2))
	}
}

// RedeemTicketV2 : redeem ticket v2
func RedeemTicketV2(c *gin.Context) {
	request := new(requests.RedeemReqV2)

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, err.Error(), "400", nil))
		logger.Warning(err.Error(), "400", false, fmt.Sprintf("%+v", request))
		return
	}

	in, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, err.Error(), "400", nil))
		logger.Warning(err.Error(), "400", false, fmt.Sprintf("%+v", request))
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	result, code, msg, msgCode, status := repositories.RedeemTicketV2(userData, request)

	c.JSON(code, builder.ApiResponseData(code, msg, msgCode, result))
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{}), string(in))
}
