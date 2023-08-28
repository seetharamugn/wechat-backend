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
		ctx.Abort()
		return
	}
	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password is required"})
		ctx.Abort()
		return
	}
	services.CreatTokenService(ctx, user)

}
