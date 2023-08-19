package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"net/http"
	"os"
)

func VerifyWebhook(ctx *gin.Context) {
	VERIFY_TOKEN := os.Getenv("VERIFY_TOKEN")
	mode := ctx.Query("hub.mode")
	token := ctx.Query("hub.verify_token")
	challenge := ctx.Query("hub.challenge")

	if mode == "subscribe" && token == VERIFY_TOKEN {
		ctx.String(http.StatusOK, challenge)
	} else {
		ctx.Status(http.StatusForbidden)
	}
}

func HandleIncomingMessage(ctx *gin.Context) {
	var messageBody Dao.WebhookMessage
	if err := ctx.ShouldBind(&messageBody); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}
	fmt.Println(messageBody)
}
