package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
)

func IncomingMessage(ctx *gin.Context, messageBody Dao.WebhookMessage) {
	if messageBody.Entry[0].Changes[0].Value.Messages[0].Type == "text" {
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Messages[0].Text.Body)
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber)
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Contacts[0].Profile.Name)
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Messages[0].ID)
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Messages[0].From)
	}
}
