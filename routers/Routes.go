package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "okay",
			})
		})
		UserRouter(v1)
		TokenRoutes(v1)
		WhatsappRouter(v1)
		TemplateRoutes(v1)
		ChatRouter(v1)
		WebhookRoutes(v1)
	}
	return router
}
