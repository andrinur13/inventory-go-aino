package middleware

import (
	"net/http"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

//Permission :
type Permission interface {
	Set(name string, callback func(*gin.Context)) gin.HandlerFunc
}

//Permit :
type Permit struct {
}

//Set :
func (p Permit) Set(name string, callback func(*gin.Context)) gin.HandlerFunc {
	if name == "PERMISSION_MASTER_USER_VIEW" || name == "PERMISSION_MASTER_USER_SAVE" {
		return callback
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusForbidden, builder.BaseResponse(false, http.StatusText(http.StatusForbidden), nil))
	}
}
