package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var ReplyUserCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "replyUser")

func IncomingMessage(ctx *gin.Context, messageBody Dao.WebhookMessage) {
	if messageBody.Entry[0].Changes[0].Value.Messages[0].Type == "text" {
		TextMessage(ctx, messageBody.Entry[0].Changes[0].Value.Messages[0].From,
			messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber,
			messageBody.Entry[0].Changes[0].Value.Messages[0].Text.Body,
			messageBody.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
			messageBody.Entry[0].Changes[0].Value.Messages[0].ID)
	}
}
func TextMessage(ctx *gin.Context, from, to, messageBody, profileName, messageId string) {
	fmt.Println(from, to, messageBody, profileName, messageId)
	var chatId interface{}
	var replyUser models.ReplyUser
	err := ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch reply user",
		})
		fmt.Println(err)
		ctx.Abort()
		return
	}
	chatId = replyUser.Id
	fmt.Println(replyUser.UserId, replyUser.Id)
	if replyUser.UserId == "" {
		userId := generateRandom()
		chatId, err = ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
		fmt.Println(chatId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to create reply user",
			})
			ctx.Abort()
			return
		}
	}

	var chat models.Chat
	err = chatCollection.FindOne(context.TODO(), bson.M{"createdBy": from}).Decode(&chat)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch chat",
		})
		ctx.Abort()
		return
	}
	if chat.CreatedBy != from {
		Numbers := []string{from, to}
		user := models.Chat{
			UserNumber:  Numbers,
			CreatedBy:   from,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			LastMessage: messageId,
		}
		_, err = chatCollection.InsertOne(context.TODO(), user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to insert chat",
				"err":   err,
			})
		}
	} else {
		_, err = chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": messageId}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to update chat",
				"err":   err,
			})
		}
	}

	message := models.Message{
		Id:            messageId,
		From:          from,
		To:            to,
		Type:          "text",
		Body:          messageBody,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReadStatus:    false,
		MessageStatus: false,
	}
	_, err = messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create template",
		})
		ctx.Abort()
		return
	}
}

func generateRandom() string {
	randNumber := 10000000 + rand.Intn(99999999-10000000)
	var user models.ReplyUser
	err := ReplyUserCollection.FindOne(context.TODO(), bson.M{"userId": randNumber}).Decode(&user)
	if err != nil {
		return strconv.Itoa(randNumber)
	}
	return generateRandom()
}
