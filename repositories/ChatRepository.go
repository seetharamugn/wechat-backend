package repositories

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

var messageCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "message")
var chatCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "chat")
var waUrl = os.Getenv("WA_URL")
var MyBucket = os.Getenv("BUCKET_NAME")
var CloudfrontUrl = os.Getenv("CLOUDFRONT_URL")

func GetAllChat(ctx *gin.Context) (interface{}, error) {
	var chats []models.Chat
	cursor, err := chatCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get templates",
			Data:       err.Error(),
		})
		ctx.Abort()
		return chats, err
	}
	for cursor.Next(context.TODO()) {
		var chat models.Chat
		err = cursor.Decode(&chat)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func SendBulkMessage(c *gin.Context, userId, templateName string, phoneNumbers []string) (interface{}, error) {
	WaAccount, err := GetAccessToken(c, userId)
	if err != nil {
		return nil, err
	}
	fmt.Println(phoneNumbers)
	for _, line := range phoneNumbers {
		payload := map[string]interface{}{
			"messaging_product": "whatsapp",
			"to":                line,
			"type":              "template",
			"template": map[string]interface{}{
				"name": templateName,
				"language": map[string]string{
					"code": "en_US",
				},
			},
		}
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		go func() {
			_, err = SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
			if err != nil {
				return
			}
		}()
	}
	return nil, nil
}

func SendTextMessage(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.TextMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "text",
		Text: models.Text{
			PreviewUrl: false,
			Body:       messageBody.MessageBody,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return nil, err
	}
	message := models.Message{
		Id:            response.Messages[0].Id,
		From:          WaAccount.PhoneNumber,
		To:            messageBody.MessageTo,
		Type:          "text",
		Body:          messageBody.MessageBody,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReadStatus:    false,
		MessageStatus: false,
	}
	resp, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create template",
		})
		ctx.Abort()
		return nil, err
	}

	return resp, nil
}

func SendReplyToTextMessage(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.TextReply{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Context: models.Context{
			MessageId: messageBody.MessageId,
		},
		Type: "text",
		Text: models.ReplyText{
			PreviewUrl: false,
			Body:       messageBody.MessageBody,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return nil, err
	}
	message := models.Message{
		Id:         response.Messages[0].Id,
		From:       WaAccount.PhoneNumber,
		To:         messageBody.MessageTo,
		Type:       "text",
		Body:       messageBody.MessageBody,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ReadStatus: false,
		ParentId:   messageBody.MessageId,
	}
	resp, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create template",
		})
		ctx.Abort()
		return nil, err
	}
	return resp, nil

}

func SendReplyByReaction(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.ReplyReaction{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "reaction",
		Reaction: models.Reaction{
			MessageId: messageBody.MessageId,
			Emoji:     messageBody.MessageBody,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func SendImageMessage(ctx *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	payload := models.ImageMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Type:             "image",
		Image: models.Image{
			Link: link,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}
	message := models.Message{
		Id:         response.Messages[0].Id,
		From:       WaAccount.PhoneNumber,
		To:         messageTo,
		Type:       "image",
		Body:       link,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ReadStatus: false,
	}
	resp, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create template",
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       resp,
	})
}

func SendVideoMessage(ctx *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		return
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		return
	}
	payload := models.VideoMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Type:             "video",
		Video: models.Video{
			Link: link,
		},
	}
	fmt.Println(payload)
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendPdfMessage(ctx *gin.Context, userId, messageTo, filename, contentType string, file multipart.File) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		return
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		return
	}
	payload := models.DocumentMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Type:             "document",
		Document: models.Document{
			Link: link,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendLocationMessage(ctx *gin.Context, messageBody models.MessageBody) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return
	}
	payload := models.LocationMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "location",
		Location: models.Location{
			Latitude:  messageBody.Latitude,
			Longitude: messageBody.Longitude,
			Name:      messageBody.LocationName,
			Address:   messageBody.LocationAddress,
		},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       response,
	})
}

func SendMessage(payload []byte, token string, phoneNumberId int, apiVersion string) (Dao.ResponseMessage, error) {
	fbUrl := waUrl + "" + apiVersion + "/" + strconv.Itoa(phoneNumberId) + "/messages"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fbUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMessage{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMessage{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMessage{}, err
	}
	var response Dao.ResponseMessage
	if res.StatusCode != http.StatusOK {
		return Dao.ResponseMessage{}, fmt.Errorf("failed to send message: %s", string(body))
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return Dao.ResponseMessage{}, fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	return response, nil
}

func UploadToS3(ctx *gin.Context, file multipart.File, filename string, contentType string) (string, error) {
	sess := initializers.ConnectAws()
	uploader := s3manager.NewUploader(sess)
	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(MyBucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": up,
		})
		return "", err
	}
	filepath := CloudfrontUrl + filename
	ctx.JSON(http.StatusOK, gin.H{
		"filepath": filepath,
	})
	return filepath, nil
}
