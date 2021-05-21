package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"twc-ota-api/db/entities"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

// AgentRouter : Routing
func AgentRouter(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	agent := r.Group("/agent")
	{
		agent.GET("/detail", permission.Set("PERMISSION_MASTER_USER_VIEW", GetAgentDetail))
		agent.POST("/update", permission.Set("PERMISSION_MASTER_USER_SAVE", UpdateAgent))
	}
	pass := r.Group("/password")
	{
		pass.POST("/update", permission.Set("PERMISSION_MASTER_USER_SAVE", UpdatePassword))
	}
	inbox := r.Group("/inbox")
	{
		inbox.GET("/list/:page/:size", permission.Set("PERMISSION_MASTER_USER_VIEW", GetInboxNotification))
	}
}

func GetAgentDetail(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.GetDetailAgent(userData)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

func UpdateAgent(c *gin.Context) {
	var param entities.AgentReq
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.UpdateProfileAgent(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

func UpdatePassword(c *gin.Context) {
	var param requests.UpdatePassword
	c.BindJSON(&param)

	in, _ := json.Marshal(param)

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	data, code, msg, stat := repositories.UpdatePassword(userData, &param)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// Inbox Notification
func GetInboxNotification(c *gin.Context) {
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

	typeNotif := c.DefaultQuery("type", "")

	// typeNotif, err := strconv.Atoi(c.Param("type"))
	// if err != nil {
	// 	c.JSON(http.StatusOK, builder.ApiResponse(false, "General Error, couldn't type size" + err.Error(), "99", nil))
	// 	logger.Warning("General Error, couldn't type size", "99", false, "")
	// 	return
	// }

	data, code, msg, stat, totalData, totalPages, currentData := repositories.GetInboxNotification(userData, typeNotif, page, size)

	out, _ := json.Marshal(data)

	contentLenght := len(string(out))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Response-Length", strconv.Itoa(contentLenght+76))
	c.Header("Transfer-Encoding", "identity")
	c.JSON(http.StatusOK, builder.ListResponse(stat, msg, code, totalData, currentData, totalPages, page, size, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}
