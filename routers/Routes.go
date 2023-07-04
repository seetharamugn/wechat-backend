package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		WhatsappRouter(v1)
		UserRouter(v1)
		TokenRoutes(v1)
	}
	return router
}
