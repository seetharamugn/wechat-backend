package services

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
	"net/http"
	"time"
)

func CreateAccount(c *gin.Context, account models.WhatsappAccount) (string, error) {
	return repositories.CreateAccount(c, account)
}

func GetAccessToken(c *gin.Context, userId string) (models.WhatsappAccount, error) {
	return repositories.GetAccessToken(c, userId)
}

func SendMsg(c *gin.Context, payload map[string]interface{}, recipient, phoneId, token string) {

	// Serialize the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error serializing JSON payload"})
		return
	}
	// Build the HTTP request
	req, err := http.NewRequest(http.MethodPost, "https://graph.facebook.com/v16.0/"+phoneId+"/messages", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error building HTTP request"})
		return
	}

	// Set the required headers
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * 10,
	}

	// Send the HTTP request and get the response
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending HTTP request"})
		return
	}
	// Close the response body to avoid resource leaks
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message to recipient " + recipient})
		return
	}
}
