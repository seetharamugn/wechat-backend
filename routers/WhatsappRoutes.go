package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/middleware"
)

func WhatsappRouter(routes *gin.RouterGroup) {
	routes.POST("/createAccount", middleware.ValidateAccessToken, controllers.CreateAccount)
}
