package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
)

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
