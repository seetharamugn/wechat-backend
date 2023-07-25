package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var VERIFY_TOKEN string

func VerifyWebhook(ctx *gin.Context) {
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
	var messageBody string
	if err := ctx.ShouldBind(&messageBody); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
		return
	}
	fmt.Println(messageBody)
}