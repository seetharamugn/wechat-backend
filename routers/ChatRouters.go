package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/middleware"
)

func ChatRouter(routes *gin.RouterGroup) {
	routes.GET("/getAllChats", middleware.ValidateAccessToken, controllers.GetAllChats)
	routes.POST("/sendBulkMessage", middleware.ValidateAccessToken, controllers.SendBulkMessage)
	routes.POST("/sendMessage", middleware.ValidateAccessToken, controllers.SendTextMessage)
	routes.POST("/sendTextMessageWithPreviewURL", middleware.ValidateAccessToken, controllers.SendTextMessageWithPreviewURL)
	routes.POST("/sendReplyByTextMessage", middleware.ValidateAccessToken, controllers.SendReplyByTextMessage)
	routes.POST("/sendReplyByReaction", middleware.ValidateAccessToken, controllers.SendReplyByReaction)
	routes.POST("/sendImage", middleware.ValidateAccessToken, controllers.SendImageMessage)
	routes.POST("/sendReplyByImageMessage", middleware.ValidateAccessToken, controllers.SendReplyByImageMessage)
	routes.POST("/sendVideo", middleware.ValidateAccessToken, controllers.SendVideoMessage)
	routes.POST("/sendReplyByVideo", middleware.ValidateAccessToken, controllers.SendReplyByVideo)
	routes.POST("/sendPdf", middleware.ValidateAccessToken, controllers.SendPdfMessage)
	routes.POST("/sendReplyByPdfMessage", middleware.ValidateAccessToken, controllers.SendReplyByPdfMessage)
	routes.POST("/sendLocation", middleware.ValidateAccessToken, controllers.SendLocationMessage)
	routes.GET("/fetchConversation", controllers.FetchConversation)
	routes.GET("/getMessagesCount", controllers.GetMessagesCount)
}
