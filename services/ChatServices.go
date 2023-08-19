package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func SendBulkMessage(c *gin.Context, userId, templateName string, phoneNumbers []string) (interface{}, error) {
	return repositories.SendBulkMessage(c, userId, templateName, phoneNumbers)
}

func SendTextMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendTextMessage(c, messageBody)

}
func SendReplyToTextMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendReplyToTextMessage(c, messageBody)
}

func SendImageMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendImageMessage(c, messageBody)
}

func SendVideoMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendVideoMessage(c, messageBody)
}
func SendPdfMessage(c *gin.Context, messageBody models.Body) (interface{}, error) {
	return repositories.SendPdfMessage(c, messageBody)
}
