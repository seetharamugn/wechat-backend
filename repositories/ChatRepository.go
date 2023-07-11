package repositories

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/seetharamugn/wachat/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

//var tokenCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "message")

var waUrl = os.Getenv("WA_URL")

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
	return response, nil
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

func SendMessage(payload []byte, token string, phoneNumberId int, apiVersion string) (interface{}, error) {
	fbUrl := waUrl + "" + apiVersion + "/" + strconv.Itoa(phoneNumberId) + "/messages"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fbUrl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var response interface{}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send message: %s", string(body))
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	return response, nil
}
