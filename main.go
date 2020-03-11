package main

import (
	"os"
	"twc-ota-api/api"
	"twc-ota-api/config"
	"twc-ota-api/db"
	"twc-ota-api/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgin"

	_ "twc-ota-api/docs"
)

// @title Dashboard API Documentation
// @version 1.0
// @description This is a documentation of using RESTfull API for Ticketing applications. json version : /swagger/doc.json

// @host 127.0.0.1:8080
// @BasePath /api/v1

func main() {
	config.Init("dev")

	log.SetLevel(log.DebugLevel)
	db.Init()
	cm := service.InitCache()

	os.Setenv("ELASTIC_APM_SERVER_URL", "http://apm-server.dev.ainosi.com/")
	defer os.Unsetenv("ELASTIC_APM_SERVER_URL")

	os.Setenv("ELASTIC_APM_SERVICE_NAME", "tnt-api")
	defer os.Unsetenv("ELASTIC_APM_SERVICE_NAME")

	os.Setenv("ELASTIC_APM_VERIFY_SERVER_CERT", "false")
	defer os.Unsetenv("ELASTIC_APM_VERIFY_SERVER_CERT")

	// router := gin.Default()
	router := gin.New()
	// router.Use(middleware.Auth(cm))
	//APM
	router.Use(apmgin.Middleware(router))

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
