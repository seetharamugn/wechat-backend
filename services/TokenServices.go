package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func CreatTokenService(ctx *gin.Context, user models.User) (interface{}, error) {
	response, err := repositories.CreateToken(ctx, user)
	return response, err
}
