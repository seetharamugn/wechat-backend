package repositories

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

	fmt.Println(messageBody)
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
			messageBody.Entry[0].Changes[0].Value.Messages[0].ID,
			messageBody.Entry[0].Changes[0].Value.Messages[0].Image.Caption)

	} else if messageBody.Entry[0].Changes[0].Value.Messages[0].Type == "video" {
		fmt.Println(messageBody.Entry[0].Changes[0].Value.Messages[0].Video)
		VideoMessage(ctx, messageBody.Entry[0].Changes[0].Value.Messages[0].From,
			messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber,
			messageBody.Entry[0].Changes[0].Value.Messages[0].Video.ID,
			messageBody.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
			messageBody.Entry[0].Changes[0].Value.Messages[0].ID,
			messageBody.Entry[0].Changes[0].Value.Messages[0].Video.Caption)

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

func ImageMessage(ctx *gin.Context, from, to, mediaId, profileName, messageId, caption string) {
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
	url, token, err := GetUrl(ctx, to, mediaId)
	if err != nil {
		return
	}
	file, err := DownLoadFile(ctx, url.Url, token, ".jpg")
	if err != nil {
		return
	}
	if chat.CreatedBy != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   to,
			LastMessage: file,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": file, "updatedAt": time.Now()}})
	}

	message := models.Message{
		Id:   messageId,
		From: from,
		To:   to,
		Type: "image",
		Body: models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "image/jpeg",
		},
		ChatId:        chatId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReadStatus:    false,
		MessageStatus: false,
	}
	messageCollection.InsertOne(context.TODO(), message)
}

func VideoMessage(ctx *gin.Context, from, to, mediaId, profileName, messageId, caption string) {
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
	url, token, err := GetUrl(ctx, to, mediaId)
	if err != nil {
		return
	}
	file, err := DownLoadFile(ctx, url.Url, token, ".mp4")
	if err != nil {
		return
	}
	if chat.CreatedBy != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   to,
			LastMessage: file,
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": from}, bson.M{"$set": bson.M{"lastMessage": file, "updatedAt": time.Now()}})
	}

	message := models.Message{
		Id:   messageId,
		From: from,
		To:   to,
		Type: "image",
		Body: models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "image/jpeg",
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

func GetUrl(c *gin.Context, phoneNumber, mediaId string) (Dao.ResponseMedia, string, error) {
	var user models.User
	fmt.Println(phoneNumber, mediaId)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNumber": phoneNumber}).Decode(&user)
	WaAccount, err := GetAccessToken(c, user.UserId)
	if err != nil {
		return Dao.ResponseMedia{}, "", err
	}
	fbUrl := waUrl + "" + WaAccount.ApiVersion + "/" + mediaId
	client := &http.Client{}
	req, err := http.NewRequest("GET", fbUrl, nil)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMedia{}, "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+WaAccount.Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMedia{}, "", err
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
		return Dao.ResponseMedia{}, "", err
	}
	fmt.Println(string(body))
	var response Dao.ResponseMedia
	if res.StatusCode != http.StatusOK {
		return Dao.ResponseMedia{}, "", fmt.Errorf("failed to send message: %s", string(body))
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return Dao.ResponseMedia{}, "", fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	fmt.Println(response.Url)
	return response, WaAccount.Token, nil
}

func DownLoadFile(ctx *gin.Context, Url string, AccessToke, fileExtension string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+AccessToke)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	file, err := UploadUrlToS3(ctx, body, fileExtension)
	if err != nil {
		return "", err
	}
	return file, nil
}

func UploadUrlToS3(ctx *gin.Context, body []byte, fileExtension string) (string, error) {
	sess := initializers.ConnectAws()
	svc := s3.New(sess)
	fileName := GenerateRandomString(12)
	// Upload the data to S3
	up, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(fileName + fileExtension),
		Body:   aws.ReadSeekCloser(bytes.NewReader(body)),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": up,
		})
		return "", err
	}
	filepath := CloudfrontUrl + fileName + fileExtension
	ctx.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})
	return filepath, nil
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}
