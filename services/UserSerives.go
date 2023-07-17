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

func GetUser(ctx *gin.Context, userId string) (models.User, error) {
	return repositories.GetUser(ctx, userId)
}

func UpdateUser(c *gin.Context, id int, body models.User) (*mongo.UpdateResult, error) {
	return repositories.UpdateUser(id, body)
}

func DeleteUser(c *gin.Context, userId string) (*mongo.DeleteResult, error) {
	return repositories.DeleteUser(c, userId)
}
