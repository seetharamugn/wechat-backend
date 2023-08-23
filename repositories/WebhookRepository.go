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
	var replyUser, chat, users map[string]interface{}
	ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	chatId = replyUser["_id"]
	fmt.Println(replyUser["userId"], replyUser["_id"])
	if replyUser["userId"] == "" {
		userId := generateRandom()
		chatId, _ = ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
		fmt.Println(chatId)
	}
	fmt.Println(replyUser, chat, users)
	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)

	fmt.Println(users["_id"])
	if chat["createdBy"] != from {
		Numbers := []interface{}{chat["_id"], users["_id"]}
		user := models.Chat{
			UserNumber:  Numbers,
			CreatedBy:   from,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			LastMessage: messageId,
		}
		chatCollection.InsertOne(context.TODO(), user)
	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": messageId}})
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
