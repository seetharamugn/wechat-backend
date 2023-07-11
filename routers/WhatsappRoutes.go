package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func WhatsappRouter(routes *gin.RouterGroup) {
	routes.POST("/createAccount", controllers.CreateAccount)

}
