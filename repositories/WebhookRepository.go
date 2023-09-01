package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
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
	} else if messageBody.Entry[0].Changes[0].Value.Messages[0].Type == "image" {
		ImageMessage(ctx, messageBody.Entry[0].Changes[0].Value.Messages[0].From,
			messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber,
			messageBody.Entry[0].Changes[0].Value.Messages[0].Image.ID,
			messageBody.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
			messageBody.Entry[0].Changes[0].Value.Messages[0].ID)

	} else if messageBody.Entry[0].Changes[0].Value.Messages[0].Type == "video" {

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
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   to,
			LastMessage: messageBody,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
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

func ImageMessage(ctx *gin.Context, from, to, mediaId, profileName, messageId string) {
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
	GetUrl(ctx, from, mediaId)
	if chat.CreatedBy != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   to,
			LastMessage: mediaId,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": mediaId, "updatedAt": time.Now()}})
	}

	message := models.Message{
		Id:   messageId,
		From: from,
		To:   to,
		Type: "text",
		Body: models.Body{
			Text: mediaId,
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

func GetUrl(c *gin.Context, phoneNumber, mediaId string) (interface{}, error) {
	var user models.User
	fmt.Println(phoneNumber, mediaId)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNumber": phoneNumber}).Decode(&user)
	WaAccount, err := GetAccessToken(c, user.UserId)
	if err != nil {
		return nil, err
	}
	fbUrl := waUrl + "" + WaAccount.ApiVersion + "/" + mediaId
	client := &http.Client{}
	req, err := http.NewRequest("GET", fbUrl, nil)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMedia{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+WaAccount.Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMedia{}, err
	}
	fmt.Println(res.StatusCode)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMedia{}, err
	}
	var response Dao.ResponseMedia
	if res.StatusCode != http.StatusOK {
		return Dao.ResponseMedia{}, fmt.Errorf("failed to send message: %s", string(body))
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return Dao.ResponseMedia{}, fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	fmt.Println(response.Url)
	return response, nil
}
