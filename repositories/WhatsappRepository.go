package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"golang.org/x/net/context"
	"net/http"
)

func CreateAccount(ctx *gin.Context, account models.WhatsappAccount) (string, error) {
	_, err := userCollection.InsertOne(context.TODO(), account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create user",
		})
		ctx.Abort()
		return "", err
	}
	return "Registration Success", nil

}
