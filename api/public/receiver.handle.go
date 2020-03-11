package public

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"twc-ota-api/middleware"

	"github.com/gin-gonic/gin"
)

// ReceiverRouter : for handling receiver json
func ReceiverRouter(r *gin.RouterGroup, permission middleware.Permission) {
	r.GET("/reader", reader)
	r.POST("/receiver", receiver)
}

func reader(c *gin.Context) {
	// c.JSON(http.StatusOK, builder.LoginResponse(true, "ok", token))
	// Open our jsonFile
	jsonFile, err := os.Open("storage/example.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened json file")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	c.JSON(http.StatusOK, result)
	// c.String(http.StatusOK, "Hahaha")
}

func receiver(c *gin.Context) {
	var request map[string]interface{}
	c.BindJSON(&request)

	file, _ := json.MarshalIndent(request, "", " ")
	_ = ioutil.WriteFile("storage/example.json", file, 0644)

	c.JSON(http.StatusOK, gin.H{"message": "Send data success", "status": "success"})
}
