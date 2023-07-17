package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func ChatRouter(routes *gin.RouterGroup) {
	routes.POST("/sendBulkMessage", controllers.SendBulkMessage)
	routes.POST("/sendMessage", controllers.SendTextMessage)
	routes.POST("/sendImage", controllers.SendImageMessage)
	routes.POST("/sendVideo", controllers.SendVideoMessage)
}
