package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/controllers"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/repositories"
	"github.com/seetharamugn/wachat/services"
)

func UserRouter(routes *gin.RouterGroup) {

	session := initializers.DBinstance()
	userRepository := repositories.NewMongoUserRepository(session)
	userService := services.NewUserService(*userRepository)
	userController := controllers.NewUserController(*userService)
	routes.POST("/signup", userController.CreateUser)
	routes.PUT("/update/:userId", userController.Update)
}
