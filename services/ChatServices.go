package services

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
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
func SendReactionMessage(c *gin.Context, userId, messageTo, body string) (interface{}, error) {
	return repositories.SendTextMessage(c, userId, messageTo, body)
}
func SendReplyByReactionMessage(c *gin.Context, userId, messageId, messageTo, body string) (interface{}, error) {
	return repositories.SendReplyByReaction(c, userId, messageId, messageTo, body)
}
func SendImageMessage(c *gin.Context, userId, messageTo, caption, link string) (interface{}, error) {
	return repositories.SendImageMessage(c, userId, messageTo, caption, link)
}

func SendReplyByImageMessage(c *gin.Context, userId, messageTo, messageId, caption, link string) (interface{}, error) {
	return repositories.SendReplyByImageMessage(c, userId, messageTo, messageId, caption, link)
}

func SendVideoMessage(c *gin.Context, userId, messageTo, caption, link string) (interface{}, error) {
	return repositories.SendVideoMessage(c, userId, messageTo, caption, link)
}

func SendReplyByVideo(c *gin.Context, userId, messageTo, messageId, caption, link string) (interface{}, error) {
	return repositories.SendReplyByVideo(c, userId, messageTo, messageId, caption, link)
}

func SendPdfMessage(c *gin.Context, userId, messageTo, caption, link string) {
	repositories.SendPdfMessage(c, userId, messageTo, caption, link)
}

func SendReplyByPdfMessage(c *gin.Context, userId, messageTo, messageId, caption, link string) (interface{}, error) {
	return repositories.SendReplyByPdfMessage(c, userId, messageTo, messageId, caption, link)
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

func UploadFile(c *gin.Context, file multipart.File, filename, contentType string) (interface{}, error) {
	return repositories.UploadToS3(c, file, filename, contentType)
}
