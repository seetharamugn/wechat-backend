package repositories

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var messageCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "message")

var waUrl = os.Getenv("WA_URL")

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

func SendTextMessage(ctx *gin.Context, messageBody models.Body) (interface{}, error) {
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

func SendReplyToTextMessage(ctx *gin.Context, messageBody models.Body) (interface{}, error) {
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

func SendImageMessage(ctx *gin.Context, messageBody models.Body) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.ImageMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "image",
		Image: models.Image{
			Link: messageBody.MessageBody,
		},
	}
	fmt.Println(payload)
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
		Type:       "image",
		Body:       messageBody.MessageBody,
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
		return nil, err
	}
	return resp, nil

}

func SendVideoMessage(ctx *gin.Context, messageBody models.Body) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.VideoMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "video",
		Video: models.Video{
			Link: messageBody.MessageBody,
		},
	}
	fmt.Println(payload)
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

func SendPdfMessage(ctx *gin.Context, messageBody models.Body) (interface{}, error) {
	WaAccount, err := GetAccessToken(ctx, strconv.Itoa(messageBody.UserId))
	if err != nil {
		return nil, err
	}
	payload := models.DocumentMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               messageBody.MessageTo,
		Type:             "document",
		Document: models.Document{
			Link: messageBody.MessageBody,
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
func SendMessage(payload []byte, token string, phoneNumberId int, apiVersion string) (Dao.ResponseMessage, error) {
	fbUrl := waUrl + "" + apiVersion + "/" + strconv.Itoa(phoneNumberId) + "/messages"
	fmt.Println(fbUrl)
	client := &http.Client{}
	req, err := http.NewRequest("POST", fbUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return Dao.ResponseMessage{}, err
	}
	fmt.Println(req.Body)
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
