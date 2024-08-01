package middleware

import (
	"context"
	"net/http"
	"time"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Create a new request with the timeout context
		c.Request = c.Request.WithContext(ctx)

		// Channel to capture the request's completion
		finished := make(chan struct{})
		panicChan := make(chan interface{})

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// Request finished normally
		case p := <-panicChan:
			// Panic occurred in the request
			panic(p)
		case <-ctx.Done():
			// Timeout occurred
			c.JSON(http.StatusGatewayTimeout, builder.ApiResponseData(http.StatusGatewayTimeout, "Request timed out", "REQUEST_TIME_OUT", nil))
			c.Abort()
		}
	}
}
