package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization") // Add Authorization header
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		UserRouter(v1)
		TokenRoutes(v1)
		WhatsappRouter(v1)
		TemplateRoutes(v1)
		ChatRouter(v1)
		WebhookRoutes(v1)
	}
	return router
}
