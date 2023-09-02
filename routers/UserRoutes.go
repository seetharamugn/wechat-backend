package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/middleware"
)

func UserRouter(routes *gin.RouterGroup) {
	routes.POST("/signup", controllers.CreateUser)
	routes.PUT("/update", middleware.ValidateAccessToken, controllers.Update)
	routes.DELETE("/delete", middleware.ValidateAccessToken, controllers.Delete)
	routes.POST("/verifyEmail", controllers.VerifyEmail)
	routes.POST("/resetPassword", controllers.ResetPassword)
	routes.GET("/getUserDetails", middleware.ValidateAccessToken, controllers.GetUserDetails)
}
