package controllers

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
)

func CreateAccount(c *gin.Context) {
	var body models.WhatsappAccount
	if body.UserId == "" || body.Token == "" || body.PhoneNumber == "" || body.PhoneNumberId == "" || body.BusinessAccountId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, _ := services.CreateAccount(c, body)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Account created successfully",
		Data:       resp,
	})
}

func SendBulkMsg(c *gin.Context) {
	phone_id := c.Query("phone_id")
	Token := c.GetHeader("Authorization")
	file, _, err := c.Request.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing CSV file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing CSV file"})
		return
	}

	// Loop through the lines and send the message to each recipient
	for _, line := range lines {
		recipient := line[0]
		// Build the request payload
		payload := map[string]interface{}{
			"messaging_product": "whatsapp",
			"to":                recipient,
			"type":              "template",
			"template": map[string]interface{}{
				"name": "hello_world",
				"language": map[string]string{
					"code": "en_US",
				},
			},
		}
		go func() {
			services.SendMsg(c, payload, recipient, phone_id, Token)
		}()

	}
	// Return a success response
	c.JSON(http.StatusOK, gin.H{"status": "Messages sent successfully"})
}

func SendMessage(c *gin.Context) {
	phoneId := c.Query("phone_id")
	Token := c.GetHeader("Authorization")
	var requestBody models.Sender

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                requestBody.To,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": false,
			"body":        requestBody.Body,
		},
	}

	go func() {
		services.SendMsg(c, payload, requestBody.To, phoneId, Token)
	}()
}

//
//func SendImage(c *gin.Context) {
//	phoneId := c.Query("phone_id")
//	Token := c.GetHeader("Authorization")
//	var requestBody models.Sender
//	payload := map[string]interface{}{
//		"messaging_product": "whatsapp",
//		"recipient_type":    "individual",
//		"to":                requestBody.To,
//		"type":              "image",
//		"image": map[string]interface{}{
//			"preview_url": false,
//			"body":        requestBody.Body,
//		},
//	}
//}

//func GetMediaId(c *gin.Context) string {
//	phoneId := c.Query("phone_id")
//	Token := c.GetHeader("Authorization")
//	media
//}
