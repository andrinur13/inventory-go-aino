package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"twc-ota-api/db/repositories"
	"twc-ota-api/logger"
	"twc-ota-api/middleware"
	"twc-ota-api/requests"
	"twc-ota-api/service"
	"twc-ota-api/utils/builder"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/v2"
)

// TicketRouter : Routing
func TicketRouterV2(r *gin.RouterGroup, permission middleware.Permission, cacheManager *service.CacheManager) {
	cm = cacheManager
	ticket := r.Group("/ticket")
	{
		ticket.POST("/redeem", permission.Set("PERMISSION_MASTER_USER_SAVE", RedeemTicketV2))
		ticket.GET("/qr", permission.Set("PERMISSION_MASTER_USER_SAVE", GetQR))
		ticket.GET("/qr/status/:qr_code", permission.Set("PERMISSION_MASTER_USER_SAVE", GetQRStatus))
		ticket.GET("/qr/summary", permission.Set("PERMISSION_MASTER_USER_SAVE", GetQRSummary))
	}
}

// RedeemTicketV2 : redeem ticket v2
func RedeemTicketV2(c *gin.Context) {
	span, spanCtx := apm.StartSpan(c.Request.Context(), "RedeemTicketV2", "handler")
	defer span.End()

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

	span.Context.SetLabel("request_body", string(in))

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	spanTime := time.Now()
	result, code, msg, msgCode, status := repositories.RedeemTicketV2(spanCtx, userData, request)
	duration := time.Since(spanTime)
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{"duration": duration}), string(in))

	c.JSON(code, builder.ApiResponseData(code, msg, msgCode, result))
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{}), string(in))
}

// GetQR : get qr
func GetQR(c *gin.Context) {
	query := new(requests.GetQrRequest)

	if err := c.ShouldBindQuery(query); err != nil {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, err.Error(), "400", nil))
		logger.Warning(err.Error(), "400", false, fmt.Sprintf("%+v", query))
		return
	}

	in, err := json.Marshal(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, err.Error(), "400", nil))
		logger.Warning(err.Error(), "400", false, fmt.Sprintf("%+v", query))
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	result, code, msg, msgCode, status := repositories.GetQRV2(userData, query)

	c.JSON(code, builder.ApiResponseData(code, msg, msgCode, result))
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{}), string(in))
}

// GetQRStatus : get qr status
func GetQRStatus(c *gin.Context) {
	qrCode := c.Param("qr_code")

	if qrCode == "" {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, "QR Code is required", "400", nil))
		logger.Warning("QR Code is required", "400", false, fmt.Sprintf("%+v", map[string]interface{}{}))
		return
	}

	in, err := json.Marshal(qrCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, builder.ApiResponse(false, err.Error(), "400", nil))
		logger.Warning(err.Error(), "400", false, fmt.Sprintf("%+v", qrCode))
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	result, code, msg, msgCode, status := repositories.GetQRStatusV2(userData, qrCode)

	c.JSON(code, builder.ApiResponseData(code, msg, msgCode, result))
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{}), string(in))
}

// GetQRSummary : get qr summary
func GetQRSummary(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	split := strings.Split(tokenString, " ")

	userData := middleware.Decode(split[1])

	result, code, msg, msgCode, status := repositories.GetQRSummaryV2(userData)

	c.JSON(code, builder.ApiResponseData(code, msg, msgCode, result))
	logger.Info(msg, strconv.Itoa(code), status, fmt.Sprintf("%v", map[string]interface{}{}), string("no request"))
}
