package main

import (
	"os"
	"time"
	"twc-ota-api/api"
	"twc-ota-api/config"
	"twc-ota-api/db"
	"twc-ota-api/middleware"
	"twc-ota-api/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgin/v2"

	_ "twc-ota-api/docs"
)

// @title Dashboard API Documentation
// @version 1.0
// @description This is a documentation of using RESTfull API for Ticketing applications. json version : /swagger/doc.json

// @host 127.0.0.1:8080
// @BasePath /api/v1

func init() {
	os.Setenv("ELASTIC_APM_SERVER_URL", "https://apm.ainosi.com")
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", "zJz71aiEaO2PRA2Wye")
	os.Setenv("ELASTIC_APM_SERVICE_NAME", "twc-api-ota")
	os.Setenv("ELASTIC_APM_ENVIRONMENT", "twc-dev")
}

func main() {
	config.Init("dev")

	log.SetLevel(log.DebugLevel)
	db.Init()
	cm := service.InitCache()

	router := gin.Default()
	// router := gin.New()
	// router.Use(middleware.Auth(cm))
	//APM
	router.Use(apmgin.Middleware(router))
	router.Use(middleware.TimeoutMiddleware(time.Duration(config.App.Timeout) * time.Second))

	api.Init(router, cm)
	api.InitWebsocket(router)

	router.Run(":" + config.App.ServerPort)
	// s := &http.Server{
	// 	Addr:           ":" + config.App.ServerPort,
	// 	Handler:        router,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.ListenAndServe()
}
