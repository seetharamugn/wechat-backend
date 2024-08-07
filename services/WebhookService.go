package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/repositories"
)

func IncomingMessage(ctx *gin.Context, messageBody Dao.WebhookResponse) {
	repositories.IncomingMessage(ctx, messageBody)

}

func WebSocketHandler(ctx *gin.Context) {
	repositories.WebSocketHandler(ctx)
}
