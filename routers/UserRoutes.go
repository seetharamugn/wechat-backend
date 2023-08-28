package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/middleware"
)

func UserRouter(routes *gin.RouterGroup) {
	routes.POST("/signup", middleware.ValidateAccessToken, controllers.CreateUser)
	routes.PUT("/update", middleware.ValidateAccessToken, controllers.Update)
	routes.DELETE("/delete", middleware.ValidateAccessToken, controllers.Delete)
	routes.POST("/verifyEmail", middleware.ValidateAccessToken, controllers.VerifyEmail)
	routes.POST("/resetPassword", middleware.ValidateAccessToken, controllers.ResetPassword)
	routes.GET("/getUserDetails", middleware.ValidateAccessToken, controllers.GetUserDetails)
}
