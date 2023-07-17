package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func CreateAccount(c *gin.Context, account models.WhatsappAccount) (string, error) {
	return repositories.CreateAccount(c, account)
}

func GetAccessToken(c *gin.Context, userId string) (models.WhatsappAccount, error) {
	return repositories.GetAccessToken(c, userId)
}
