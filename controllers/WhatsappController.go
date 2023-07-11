package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
)

func CreateAccount(c *gin.Context) {
	var body models.WhatsappAccount
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if body.PhoneNumber == "" || body.UserId == 0 || body.PhoneNumberId == 0 || body.BusinessAccountId == 0 || body.Token == "" || body.ApiVersion == "" {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "All fields are required",
			Data:       nil,
		})
		c.Abort()
		return
	}
	resp, _ := services.CreateAccount(c, body)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Account created successfully",
		Data:       resp,
	})
}

func GetAccessToken(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "userId is required",
			Data:       nil,
		})
		c.Abort()
		return
	}
	resp, _ := services.GetAccessToken(c, userId)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Access token fetched successfully",
		Data:       resp,
	})
	return
}
