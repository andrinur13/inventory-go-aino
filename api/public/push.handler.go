package public

import (
	"log"
	"net/http"
	"twc-ota-api/middleware"
	"twc-ota-api/utils/builder"

	"github.com/NaySoftware/go-fcm"
	"github.com/gin-gonic/gin"
)

const (
	serverKey = "AAAAHYh8K-k:APA91bFfHItpO6P3KKfY3-kt6CZFDtBMmtgLK7TrNn8z3Ew5asJaAOUWhQW54kVsMqcXVg4p3C5fpid_48-8m6-O1ytqOPyABULLRyjj0H4SVR-zKxLx06G4UO3IBMwOUIo89U_66rMT"
	topic     = "/topics/pushNotif"
)

type Send struct {
	To           string `json:"to"`
	Notification Notif  `json:"notification"`
}

type Notif struct {
	Title       string      `json:"title"`
	Body        interface{} `json:"body"`
	Icon        string      `json:"icon"`
	ClickAction string      `json:"click_action"`
}

// PushRouter : for handling push to fcm
func PushRouter(r *gin.RouterGroup, permission middleware.Permission) {
	r.POST("/push", push)
}

func push(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)

	log.Print(data["fcmToken"])

	token := data["fcmToken"].(string)

	ids := []string{
		token,
	}

	notif := Notif{
		Title:       "Title",
		Body:        data["data"],
		Icon:        "https://example.com/icon.png",
		ClickAction: "Post Link",
	}

	send := Send{
		To:           token,
		Notification: notif,
	}
	// send := map[string]map[string]string{
	// 	"to": token,
	// 	"notification": map[string]string{
	// 		"title":        "Title",
	// 		"body":         data["data"],
	// 		"icon":         "https://example.com/icon.png",
	// 		"click_action": "Post Link",
	// 	},
	// }

	f := fcm.NewFcmClient(serverKey)
	// f.NewFcmMsgTo(topic, data)
	f.NewFcmRegIdsMsg(ids, send)

	status, err := f.Send()

	if err == nil {
		// status.PrintResults()
		// resp := data["data"]
		c.JSON(http.StatusOK, builder.ApiResponse(true, "Send push notification data success", "01", status))
	} else {
		c.JSON(http.StatusOK, builder.ApiResponse(false, "Error occured: "+err.Error(), "99", nil))
	}
}
