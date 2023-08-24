package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"mime/multipart"
)

func GetAllChat(c *gin.Context) (interface{}, error) {
	return repositories.GetAllChat(c)
}

func SendBulkMessage(c *gin.Context, userId, templateName string, phoneNumbers []string) (interface{}, error) {
	return repositories.SendBulkMessage(c, userId, templateName, phoneNumbers)
}

func SendTextMessage(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendTextMessage(c, messageBody)

}
func SendReplyToTextMessage(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendReplyToTextMessage(c, messageBody)
}

func SendReplyByReaction(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendReplyByReaction(c, messageBody)
}
func SendImageMessage(c *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	repositories.SendImageMessage(c, userId, messageTo, filename, contentType, file)
}

func SendVideoMessage(c *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	repositories.SendVideoMessage(c, userId, messageTo, filename, contentType, file)
}
func SendPdfMessage(c *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	repositories.SendPdfMessage(c, userId, messageTo, filename, contentType, file)
}

func SendLocationMessage(c *gin.Context, messageBody models.MessageBody) {
	repositories.SendLocationMessage(c, messageBody)
}
