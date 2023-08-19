package main

import (
	"github.com/gin-gonic/gin"
	initializers "github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/routers"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	route := routers.SetupRouter()
	route.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "okay",
		})
	})
	route.Run()
}
