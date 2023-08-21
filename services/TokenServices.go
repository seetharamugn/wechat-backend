package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func CreatTokenService(ctx *gin.Context, user models.User) {
	repositories.CreateToken(ctx, user)

}
