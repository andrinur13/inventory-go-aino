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
