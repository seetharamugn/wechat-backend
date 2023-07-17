package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func UserRouter(routes *gin.RouterGroup) {
	routes.POST("/signup", controllers.CreateUser)
	routes.GET("/getUser", controllers.GetUser)
	routes.PUT("/update/:userId", controllers.Update)
	routes.DELETE("/delete/:userId", controllers.Delete)
}
