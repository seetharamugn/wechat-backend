package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("", func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "working"})
	})
	v1 := router.Group("/api/v1")
	{
		WhatsappRouter(v1)
		UserRouter(v1)
		TokenRoutes(v1)
	}
	return router
}
