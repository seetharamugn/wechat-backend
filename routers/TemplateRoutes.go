package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
)

func TemplateRoutes(routes *gin.RouterGroup) {
	routes.POST("/createTemplate", controllers.CreateTemplate)
	routes.GET("/getTemplate", controllers.GetTemplate)
	routes.PUT("/updateTemplate", controllers.UpdateTemplate)
	routes.DELETE("/deleteTemplate", controllers.DeleteTemplate)
	routes.GET("/getTemplateList", controllers.GetTemplateList)
}
