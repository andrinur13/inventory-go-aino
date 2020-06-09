package middleware

import (
	"net/http"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

//Auth : middleware for auth JWT
func Auth(cache *service.CacheManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		_, err := Authorize(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, builder.ApiResponse(false, err.Error(), "99", nil))
			c.Abort()
			return
			// log.Print(err)
		}
		c.Next()
	}
}
