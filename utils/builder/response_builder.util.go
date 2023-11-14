package builder

// response message
const (
	MessageFetchTrxFailed  = "Failed to get transactions data"
	MessageFetchTrxSuccess = "Transactions data retrieved successfully"
	MessageAuthFailed      = "Authentication failed"
)

type ErrResponse struct {
	Error string
}

type MsgResponse struct {
	Message string
}

type ReadyStatement struct {
	Dbconn int
	Query  string
	Params []interface{}
}

type ResponseData struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	MessageCode string      `json:"message_code"`
	Data        interface{} `json:"data"`
}

func BaseResponse(success bool, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data}
}

func LoginResponse(success bool, message string, code string, token string) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"message": message,
		"code":    code,
		"token":   token,
		"type":    "Bearer",
	}
}

func ApiResponse(success bool, message string, code string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"message": message,
		"code":    code,
		"data":    data,
	}
}

// func PushResponse(success bool, message string, code string, data interface{}, notifData interface{}) map[string]interface{} {
// 	return map[string]interface{}{
// 		"success":      success,
// 		"message":      message,
// 		"code":         code,
// 		"data":         data,
// 		"fcm_response": notifData,
// 	}
// }

func ListResponse(success bool, message string, code string, countData, currentData, totalPages, page, dataLimit, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success":           success,
		"message":           message,
		"code":              code,
		"data_total":        countData,
		"data_current_page": currentData,
		"page_total":        totalPages,
		"page_current":      page,
		"data_limit":        dataLimit,
		"data":              data,
	}
}

func WebsocketResponse(success bool, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data}
}

func ApiResponseData(code int, message string, messageCode string, data interface{}) *ResponseData {
	return &ResponseData{
		Code:        code,
		Message:     message,
		MessageCode: messageCode,
		Data:        data,
	}
}
