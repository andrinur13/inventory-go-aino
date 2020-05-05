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

// DiscountRouter : Routing Discount
func DiscountRouter(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	discount := r.Group("/discount")
	{
		discount.GET("/agent", permission.Set("PERMISSION_MASTER_USER_VIEW", GetDiscountAgent))
		discount.GET("/multidestination", permission.Set("PERMISSION_MASTER_USER_VIEW", GetDiscountDest))
	}
}

// GetDiscountAgent : get discount agent
func GetDiscountAgent(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetDiscountMulti(userData, "AGENT")

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

// GetDiscountDest : get discount multidestination
func GetDiscountDest(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetDiscountMulti(userData, "MULTIDESTINATION")

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}
