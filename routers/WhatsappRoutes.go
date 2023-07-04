package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func WhatsappRouter(routes *gin.RouterGroup) {
	routes.POST("/send-bulk-message", controllers.SendBulkMsg)
	routes.POST("/send-message", controllers.SendMessage)
}
