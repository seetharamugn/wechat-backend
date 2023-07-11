package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func SendTextMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendTextMessage(c, messageBody)

}

func SendImageMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendImageMessage(c, messageBody)
}
