package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
)

func GetAllChats(ctx *gin.Context) {
	PhoneNumber := ctx.Query("phoneNumber")
	if PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "phoneNumber is required",
			Data:       nil,
		})
		return
	}
	resp, err := services.GetAllChat(ctx, PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       resp,
	})
}
func SendBulkMessage(c *gin.Context) {
	userId := c.PostForm("userId")
	templateName := c.PostForm("templateName")
	file, _, err := c.Request.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing CSV file"})
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var contacts []string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing CSV file"})
			return
		}
		if len(row) > 1 {
			contacts = append(contacts, row[1])
		}
	}
	if len(contacts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No contact data found in the CSV"})
		return
	}
	fmt.Println(contacts)
	response, _ := services.SendBulkMessage(c, userId, templateName, contacts[1:])
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendTextMessage(c *gin.Context) {
	var requestBody models.MessageBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}

	userId := requestBody.UserId
	messageTo := requestBody.MessageTo
	body := requestBody.MessageBody
	messageId := requestBody.MessageId
	messageType := requestBody.MessageType
	link := requestBody.File.Url

	switch messageType {
	case "image":
		if messageId != "" {
			services.SendReplyByImageMessage(c, userId, messageTo, messageId, body, link)
		} else {
			services.SendImageMessage(c, userId, messageTo, body, link)
		}
	case "video":
		if messageId != "" {
			services.SendReplyByVideo(c, userId, messageTo, messageId, body, link)
		} else {
			services.SendVideoMessage(c, userId, messageTo, body, link)
		}
	case "audio":
	case "document":
		if messageId != "" {
			services.SendReplyByPdfMessage(c, userId, messageTo, messageId, body, link)
		} else {
			services.SendPdfMessage(c, userId, messageTo, body, link)
		}
	case "reaction":
		if messageId != "" {
			services.SendReplyByReaction(c, userId, messageId, messageTo, body)
		}
	default:
		if messageId != "" {
			services.SendReplyByTextMessage(c, userId, messageId, messageTo, body)
		} else {
			services.SendTextMessage(c, userId, messageTo, body)
		}
	}

}

func SendLocationMessage(c *gin.Context) {
	var requestBody models.MessageBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := services.SendLocationMessage(c, requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
	} else {
		c.JSON(http.StatusOK, Dao.Response{
			StatusCode: http.StatusOK,
			Message:    "Message sent successfully",
			Data:       response,
		})
	}
}

func FetchConversation(c *gin.Context) {
	chatId := c.Query("chatId")
	response, err := services.FetchConversation(c, chatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "get Message successfully",
		Data:       response,
	})

}

func GetMessagesCount(c *gin.Context) {
	phoneNumber := c.Query("PhoneNumber")
	response, err := services.GetMessagesCount(c, phoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message count fetch successfully",
		Data:       response,
	})
}

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var contentType, filename string
	if file != nil {
		contentType = header.Header.Get("Content-Type")
		filename = header.Filename
	}
	services.UploadFile(c, file, filename, contentType)
}
