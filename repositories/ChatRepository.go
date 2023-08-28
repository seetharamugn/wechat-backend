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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
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
var Chat models.Chat

func GetAllChat(ctx *gin.Context, PhoneNumber string) (interface{}, error) {
	var chats []models.Chat
	cursor, err := chatCollection.Find(context.TODO(), bson.M{"createdBy": PhoneNumber})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get chats",
		})
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var chat models.Chat
		err = cursor.Decode(&chat)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode chat",
			})
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
	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}, bson.M{"$set": bson.M{"lastMessage": messageBody.MessageBody}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageBody.MessageTo, messageBody.MessageBody, "", "", "text")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendTextMessageWithPreviewURL(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.PreviewUrl{
		MessagingProduct: "whatsapp",
		To:               messageBody.MessageTo,
		Text: models.Text{
			PreviewUrl: true,
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
	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}, bson.M{"$set": bson.M{"lastMessage": messageBody.MessageBody}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageBody.MessageTo, messageBody.MessageBody, "", "", "text")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendReplyByTextMessage(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}, bson.M{"$set": bson.M{"lastMessage": messageBody.MessageBody}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageBody.MessageTo, messageBody.MessageBody, "", messageBody.MessageId, "text")
	if err != nil {
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageBody.MessageTo}, bson.M{"$set": bson.M{"lastMessage": messageBody.MessageBody}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageBody.MessageTo, messageBody.MessageBody, "", messageBody.MessageId, "reaction")
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, "", link, "", "image")
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       resp,
	})
}

func SendReplyByImageMessage(ctx *gin.Context, userId, messageTo, messageId, filename, contentType string, file multipart.File) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		return nil, err
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return nil, err
	}
	payload := models.ImageReply{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Context: models.Context{
			MessageId: messageId,
		},
		Type: "image",
		Image: models.Image{
			Link: link,
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, "", link, messageId, "image")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendVideoMessage(ctx *gin.Context, userId, messageTo, caption, filename, contentType string, file multipart.File) {
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
			Link:    link,
			Caption: caption,
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, caption, link, "", "video")
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       resp,
	})
}

func SendReplyByVideo(ctx *gin.Context, userId, messageTo, messageId, caption, filename, contentType string, file multipart.File) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		return nil, err
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return nil, err
	}
	payload := models.VideoReply{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Context: models.Context{
			MessageId: messageId,
		},
		Type: "video",
		Video: models.Video{
			Link:    link,
			Caption: caption,
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, caption, link, messageId, "video")
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, "", link, "", "document")
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Message sent successfully",
		Data:       resp,
	})
}

func SendReplyByPdfMessage(ctx *gin.Context, userId, messageTo, messageId, caption, filename, contentType string, file multipart.File) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, userId)
	if err != nil {
		return nil, err
	}
	link, err := UploadToS3(ctx, file, filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return nil, err
	}
	payload := models.DocumentReply{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageTo,
		Context: models.Context{
			MessageId: messageId,
		},
		Type: "document",
		Document: models.Document{
			Link:    link,
			Caption: caption,
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

	chatCollection.FindOne(context.TODO(), bson.M{"createdBy": messageTo}).Decode(&Chat)
	chatCollection.UpdateOne(context.TODO(), bson.M{"createdBy": messageTo}, bson.M{"$set": bson.M{"lastMessage": link}})
	resp, err := InsertMessageIntoDB(ctx, Chat.ID, response.Messages[0].Id, WaAccount.PhoneNumber, messageTo, caption, link, messageId, "document")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func SendLocationMessage(ctx *gin.Context, messageBody models.MessageBody) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	response, err := SendMessage(jsonBody, WaAccount.Token, WaAccount.PhoneNumberId, WaAccount.ApiVersion)
	if err != nil {
		return nil, err
	}
	return response, err
}

func FetchConversation(ctx *gin.Context, chatId string) ([]models.Message, error) {
	objectId, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		log.Println("Invalid id")
	}
	var messages []models.Message
	options := options.Find()
	options.SetSort(bson.M{"timestamp": -1}) // Sort by timestamp in descending order
	options.SetLimit(20)                     // Limit to 20 messages

	cursor, err := messageCollection.Find(context.TODO(), bson.M{"chatId": objectId}, options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get messages",
		})
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var message models.Message
		err = cursor.Decode(&message)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode message",
			})
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
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

func InsertMessageIntoDB(ctx *gin.Context, chatId primitive.ObjectID, messageId, phoneNumber, messageTo, messageBody, link, parentMessageId, messageType string) (*mongo.InsertOneResult, error) {
	Body := models.Body{}
	if messageType == "text" {
		Body = models.Body{
			Text: messageBody,
		}
	} else if messageType == "image" {
		Body = models.Body{
			Text:     messageBody,
			Url:      link,
			MimeType: "image/jpeg",
		}
	} else if messageType == "video" {
		Body = models.Body{
			Text:     messageBody,
			Url:      link,
			MimeType: "video/mp4",
		}
	} else if messageType == "document" {
		Body = models.Body{
			Text:     messageBody,
			Url:      link,
			MimeType: "application/pdf",
		}
	} else if messageType == "reaction" {
		Body = models.Body{
			Text: messageBody,
		}
	}
	message := models.Message{
		Id:            messageId,
		From:          phoneNumber,
		To:            messageTo,
		Type:          messageType,
		Body:          Body,
		ChatId:        chatId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReadStatus:    false,
		MessageStatus: false,
		ParentId:      parentMessageId,
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
