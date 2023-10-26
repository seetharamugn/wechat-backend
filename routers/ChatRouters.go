package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func ChatRouter(routes *gin.RouterGroup) {
	routes.GET("/getAllChats", controllers.GetAllChats)
	routes.POST("/sendBulkMessage", controllers.SendBulkMessage)
	routes.POST("/sendMessage", controllers.SendTextMessage)
	routes.POST("/sendTextMessageWithPreviewURL", controllers.SendTextMessageWithPreviewURL)
	routes.POST("/sendReplyByTextMessage", controllers.SendReplyByTextMessage)
	routes.POST("/sendReplyByReaction", controllers.SendReplyByReaction)
	routes.POST("/sendImage", controllers.SendImageMessage)
	routes.POST("/sendReplyByImageMessage", controllers.SendReplyByImageMessage)
	routes.POST("/sendVideo", controllers.SendVideoMessage)
	routes.POST("/sendReplyByVideo", controllers.SendReplyByVideo)
	routes.POST("/sendPdf", controllers.SendPdfMessage)
	routes.POST("/sendReplyByPdfMessage", controllers.SendReplyByPdfMessage)
	routes.POST("/sendLocation", controllers.SendLocationMessage)
	routes.GET("/fetchConversation", controllers.FetchConversation)
	routes.GET("/getMessagesCount", controllers.GetMessagesCount)
}
