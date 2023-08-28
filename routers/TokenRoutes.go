package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func TokenRoutes(routes *gin.RouterGroup) {
	routes.POST("/login", controllers.TokenControllers)
	//routes.POST("/refresh", controller.RefreshToken)
}
