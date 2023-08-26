package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(ctx *gin.Context, user models.User) {
	repositories.CreateUser(ctx, user)
}

func UpdateUser(ctx *gin.Context, userId string, body models.User) (*mongo.UpdateResult, error) {
	return repositories.UpdateUser(ctx, userId, body)
}

func DeleteUser(c *gin.Context, userId int) {
	repositories.DeleteUser(c, userId)
}

func VerifyEmail(c *gin.Context, email string) {
	repositories.VerifyEmail(c, email)
}

func ResetPassword(c *gin.Context, email, password string) {
	repositories.ResetPassword(c, email, password)
}

func GetUserDetails(c *gin.Context, userId string) {
	repositories.GetUserDetails(c, userId)

}
