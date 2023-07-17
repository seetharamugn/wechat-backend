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
	var requestBody models.Body

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

func SendImageMessage(c *gin.Context) {

	var requestBody models.Body
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}
	response, _ := services.SendImageMessage(c, requestBody)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendVideoMessage(c *gin.Context) {
	var requestBody models.Body
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil})
		return
	}
	response, _ := services.SendVideoMessage(c, requestBody)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}
