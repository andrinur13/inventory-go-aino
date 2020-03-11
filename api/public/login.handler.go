package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

// LoginRouter : for handling login route
func LoginRouter(r *gin.RouterGroup, permission middleware.Permission) {
	r.POST("/login", login)
}

func login(c *gin.Context) {
	// log.Print(config.App.SvdHost)
	var req interface{}
	c.BindJSON(&req)
	in, _ := json.Marshal(req)
	// log.Print(dataMerchant)

	// token, _ := middleware.CreateJwtToken(dataMerchant)
	data, code, msg, stat := repositories.GetUser(req)

	// c.JSON(http.StatusOK, builder.BaseResponse(true, "ok", nil))
	c.JSON(http.StatusOK, builder.ApiResponse(stat, msg, code, data))
	logger.Info(msg, code, stat, fmt.Sprintf("%v", data), string(in))
}
