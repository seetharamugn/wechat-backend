package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func WebhookRoutes(routes *gin.RouterGroup) {
	routes.GET("/ws", controllers.WebSocketHandler)
	routes.GET("/webhook", controllers.VerifyWebhook)
	routes.POST("/webhook", controllers.HandleIncomingMessage)

}
