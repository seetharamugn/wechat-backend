package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
)

func TokenControllers(ctx *gin.Context) {
	var user models.User
	if ctx.Bind(&user) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read the input",
		})
		return
	}
	if user.Username == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		ctx.Abort()
		return
	}
	response, _ := services.CreatTokenService(ctx, user)
	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data":       response,
	})
}
