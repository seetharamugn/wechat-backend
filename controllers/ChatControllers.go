package controllers

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"io"
	"net/http"
)

func GetAllChats(ctx *gin.Context) {
	resp, _ := services.GetAllChat(ctx)
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

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, _ := services.SendTextMessage(c, requestBody)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendReplyToTextMessage(c *gin.Context) {
	var requestBody models.MessageBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}
	response, _ := services.SendReplyToTextMessage(c, requestBody)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendReplyByReaction(c *gin.Context) {
	var requestBody models.MessageBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}
	response, err := services.SendReplyByReaction(c, requestBody)
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
func SendImageMessage(c *gin.Context) {
	userId := c.PostForm("userId")
	messageTo := c.PostForm("messageTo")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
	}
	filename := header.Filename
	contentType := header.Header.Get("Content-Type")

	services.SendImageMessage(c, userId, messageTo, filename, contentType, file)
}

func SendVideoMessage(c *gin.Context) {
	userId := c.PostForm("userId")
	messageTo := c.PostForm("messageTo")
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
	}
	filename := header.Filename
	contentType := header.Header.Get("Content-Type")

	services.SendVideoMessage(c, userId, messageTo, filename, contentType, file)
}

func SendPdfMessage(c *gin.Context) {
	userId := c.PostForm("userId")
	messageTo := c.PostForm("messageTo")
	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
	}
	filename := header.Filename
	contentType := header.Header.Get("Content-Type")

	services.SendPdfMessage(c, userId, messageTo, filename, contentType, file)
}

func SendLocationMessage(c *gin.Context) {
	var requestBody models.MessageBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.SendLocationMessage(c, requestBody)
}
