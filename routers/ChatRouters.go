package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func ChatRouter(routes *gin.RouterGroup) {
	routes.GET("/getAllChats", controllers.GetAllChats)
	routes.POST("/sendBulkMessage", controllers.SendBulkMessage)
	routes.POST("/sendMessage", controllers.SendTextMessage)
	routes.POST("/sendReplyToTextMessage", controllers.SendReplyToTextMessage)
	routes.POST("/sendReplyByReaction", controllers.SendReplyByReaction)
	routes.POST("/sendImage", controllers.SendImageMessage)
	//routes.POST("/sendReplyToImageMessage", controllers.SendReplyToImageMessage)
	routes.POST("/sendVideo", controllers.SendVideoMessage)
	routes.POST("/sendPdf", controllers.SendPdfMessage)
	routes.POST("/sendLocation", controllers.SendLocationMessage)
}
