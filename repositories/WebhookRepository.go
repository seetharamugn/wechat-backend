package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"math/rand"
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
	var chatId interface{}
	var replyUser models.ReplyUser
	var chat models.Chat
	var users models.User
	ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	chatId = replyUser.ID
	if replyUser.UserId == "" {
		userId := generateRandom()
		ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
	}
	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)
	chatId = chat.ID
	if chat.CreatedBy != from {
		Numbers := []interface{}{replyUser.ID, users.ID}
		user := models.Chat{
			UserNumber:  Numbers,
			CreatedBy:   to,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			LastMessage: messageBody,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": messageBody, "updatedAt": time.Now()}})
	}

	message := models.Message{
		Id:   messageId,
		From: from,
		To:   to,
		Type: "text",
		Body: models.Body{
			Text: messageBody,
		},
		ChatId:        chatId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReadStatus:    false,
		MessageStatus: false,
	}
	messageCollection.InsertOne(context.TODO(), message)
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
