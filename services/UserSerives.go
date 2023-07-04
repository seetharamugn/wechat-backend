package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

type UserService struct {
	mongoRepository repositories.MongoUserRepository
}

func NewUserService(repository repositories.MongoUserRepository) *UserService {
	return &UserService{mongoRepository: repository}
}

func (service *UserService) CreateUser(ctx *gin.Context, user models.User) (string, error) {
	return service.mongoRepository.CreateUser(ctx, user)
}
