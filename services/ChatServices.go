package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"mime/multipart"
)

func GetAllChat(c *gin.Context, PhoneNumber string) (interface{}, error) {
	return repositories.GetAllChat(c, PhoneNumber)
}

func SendBulkMessage(c *gin.Context, userId, templateName string, phoneNumbers []string) (interface{}, error) {
	return repositories.SendBulkMessage(c, userId, templateName, phoneNumbers)
}

func SendTextMessage(c *gin.Context, userId, messageTo, body string) {
	repositories.SendTextMessage(c, userId, messageTo, body)

}

func SendTextMessageWithPreviewURL(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendTextMessageWithPreviewURL(c, messageBody)
}

func SendReplyByTextMessage(c *gin.Context, userId, messageId, messageTo, body string) (interface{}, error) {
	return repositories.SendReplyByTextMessage(c, userId, messageId, messageTo, body)
}

func SendReplyByReaction(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendReplyByReaction(c, messageBody)
}
func SendImageMessage(c *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	repositories.SendImageMessage(c, userId, messageTo, filename, contentType, file)
}

func SendReplyByImageMessage(c *gin.Context, userId, messageTo, messageId, filename, contentType string, file multipart.File) (interface{}, error) {
	return repositories.SendReplyByImageMessage(c, userId, messageTo, messageId, filename, contentType, file)
}

func SendVideoMessage(c *gin.Context, userId, messageTo, caption, filename, contentType string, file multipart.File) {
	repositories.SendVideoMessage(c, userId, messageTo, caption, filename, contentType, file)
}

func SendReplyByVideo(c *gin.Context, userId, messageTo, messageId, caption, filename, contentType string, file multipart.File) (interface{}, error) {
	return repositories.SendReplyByVideo(c, userId, messageTo, messageId, caption, filename, contentType, file)
}

func SendPdfMessage(c *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	repositories.SendPdfMessage(c, userId, messageTo, filename, contentType, file)
}

func SendReplyByPdfMessage(c *gin.Context, userId, messageTo, messageId, caption, filename, contentType string, file multipart.File) (interface{}, error) {
	return repositories.SendReplyByPdfMessage(c, userId, messageTo, messageId, caption, filename, contentType, file)
}

func SendLocationMessage(c *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	return repositories.SendLocationMessage(c, messageBody)
}

func FetchConversation(c *gin.Context, chatId string) (interface{}, error) {
	return repositories.FetchConversation(c, chatId)
}

func GetMessagesCount(c *gin.Context, phoneNumber string) (interface{}, error) {
	return repositories.GetMessagesCount(c, phoneNumber)
}
