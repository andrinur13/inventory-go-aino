package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twc-ota-api/db/entities"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

// PublicRouter : for handling public
func PublicRouter(r *gin.RouterGroup, permission middleware.Permission) {
	r.GET("/agent", agent)
	r.POST("/create/agent", RegisterAgent)
	r.POST("/password/reset", ResetPassword)
	r.POST("/password/update", UpdatePassword)
}

func agent(c *gin.Context) {
	data, code, msg, stat := repositories.GetAgent()

	// c.JSON(http.StatusOK, builder.BaseResponse(true, "ok", nil))
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

// RegisterAgent : register new agent
func RegisterAgent(c *gin.Context) {
	var req entities.AgentReq
	c.BindJSON(&req)
	in, _ := json.Marshal(req)

	data, code, msg, stat := repositories.InsertAgent(nil, &req)

	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// ResetPassword : reset password
func ResetPassword(c *gin.Context) {
	var req requests.ResetPassword
	c.BindJSON(&req)
	in, _ := json.Marshal(req)

	data, code, msg, stat := repositories.ResetPassword(&req)

	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}

// UpdatePassword : for update resetted password
func UpdatePassword(c *gin.Context) {
	var req requests.UpdateResetPassword
	c.BindJSON(&req)
	in, _ := json.Marshal(req)

	data, code, msg, stat := repositories.UpdateResetPassword(&req)

	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}
