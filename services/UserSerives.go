package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(ctx *gin.Context, user models.User) (*mongo.InsertOneResult, error) {
	return repositories.CreateUser(ctx, user)
}

func GetUser(ctx *gin.Context, userId int) (models.User, error) {
	return repositories.GetUser(ctx, userId)
}

func UpdateUser(ctx *gin.Context, id int, body models.User) (*mongo.UpdateResult, error) {
	return repositories.UpdateUser(ctx, id, body)
}

func DeleteUser(c *gin.Context, userId int) (*mongo.DeleteResult, error) {
	return repositories.DeleteUser(c, userId)
}
