package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func UserRouter(routes *gin.RouterGroup) {
	routes.POST("/signup", controllers.CreateUser)
	routes.PUT("/update", controllers.Update)
	routes.DELETE("/delete", controllers.Delete)
	routes.POST("/verifyEmail", controllers.VerifyEmail)
	routes.POST("/resetPassword", controllers.ResetPassword)
	routes.GET("/getUserDetails", controllers.GetUserDetails)
}
