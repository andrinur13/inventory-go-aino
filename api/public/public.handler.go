package public

import (
	"fmt"
	"net/http"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/utils/builder"
	"twc-ota-api/db/entities"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// PublicRouter : for handling public
func PublicRouter(r *gin.RouterGroup, permission middleware.Permission) {
	r.GET("/agent", agent)
	r.POST("/create/agent", RegisterAgent)
}

func agent(c *gin.Context) {
	data, code, msg, stat := repositories.GetAgent()

	// c.JSON(http.StatusOK, builder.BaseResponse(true, "ok", nil))
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), "")
}

func RegisterAgent(c *gin.Context) {
	var req entities.AgentReq
	c.BindJSON(&req)
	in, _ := json.Marshal(req)

	data, code, msg, stat := repositories.InsertAgent(nil, &req)

	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}
