package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func ChatRouter(routes *gin.RouterGroup) {
	routes.GET("/getAllChats", controllers.GetAllChats)
	routes.POST("/sendBulkMessage", controllers.SendBulkMessage)
	routes.POST("/sendMessage", controllers.SendTextMessage)
	routes.POST("/sendLocation", controllers.SendLocationMessage)
	routes.GET("/fetchConversation", controllers.FetchConversation)
	routes.GET("/getMessagesCount", controllers.GetMessagesCount)
}
