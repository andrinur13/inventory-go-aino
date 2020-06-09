package api

import (
	"twc-ota-api/api/master"
	"twc-ota-api/api/public"
	"twc-ota-api/middleware"
	"twc-ota-api/service"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//Init : registering router handler
func Init(router *gin.Engine, cache *service.CacheManager) {
	permission := middleware.Permit{}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(public.NotFound)

	auth := router.Group("/auth")
	{
		public.LoginRouter(auth, permission)
	}

	cm := service.InitCache()

	route := router.Group("/api")
	route.Use(middleware.Auth(cm))
	{
		// master.UserRouter(v1, permission, cache)
		master.TicketRouter(route, permission, cache)
		master.DiscountRouter(route, permission, cache)
	}

	pub := router.Group("/public")
	{
		public.PublicRouter(pub, permission)
	}
}
