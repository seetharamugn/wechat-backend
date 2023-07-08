package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	mongoRepository repositories.MongoUserRepository
}

func NewUserService(repository repositories.MongoUserRepository) *UserService {
	return &UserService{mongoRepository: repository}
}

func (service *UserService) CreateUser(ctx *gin.Context, user models.User) (*mongo.InsertOneResult, error) {
	return service.mongoRepository.CreateUser(ctx, user)
}

func (service *UserService) UpdateUser(c *gin.Context, id int, body models.User) (*mongo.UpdateResult, error) {
	return service.mongoRepository.UpdateUser(id, body)
}
