package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/middleware"
)

func TemplateRoutes(routes *gin.RouterGroup) {
	routes.POST("/createTemplate", middleware.ValidateAccessToken, controllers.CreateTemplate)
	routes.GET("/getTemplate", middleware.ValidateAccessToken, controllers.GetTemplate)
	routes.PUT("/updateTemplate", middleware.ValidateAccessToken, controllers.UpdateTemplate)
	routes.DELETE("/deleteTemplate", middleware.ValidateAccessToken, controllers.DeleteTemplate)
	routes.GET("/getTemplateList", middleware.ValidateAccessToken, controllers.GetTemplateList)
}
