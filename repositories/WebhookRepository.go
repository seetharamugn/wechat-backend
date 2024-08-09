package repositories

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var ReplyUserCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "replyUser")

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections (adjust for production)
		},
	}
	clients   = make(map[string]*websocket.Conn) // Connected clients
	clientsMu sync.Mutex                         // Mutex for safe concurrent access
)

func IncomingMessage(ctx *gin.Context, messageBody Dao.WebhookResponse) {

	fmt.Println(messageBody)
	var messageBodyText, to, from, phoneNumber, profileName, msgID, messageType, messageStatusID, messageStatus, recipentId, emoji, id string
	// Assuming messageBody is of type WebhookResponse
	if len(messageBody.Entry) > 0 &&
		len(messageBody.Entry[0].Changes) > 0 &&
		len(messageBody.Entry[0].Changes[0].Value.Messages) > 0 {

		// Access the message body
		//check the Text
		to = messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber
		from = messageBody.Entry[0].Changes[0].Value.Messages[0].From
		phoneNumber = messageBody.Entry[0].Changes[0].Value.Metadata.DisplayPhoneNumber
		profileName = messageBody.Entry[0].Changes[0].Value.Contacts[0].Profile.Name
		msgID = messageBody.Entry[0].Changes[0].Value.Messages[0].ID
		messageType = messageBody.Entry[0].Changes[0].Value.Messages[0].Type
		if messageType == "text" {
			messageBodyText = messageBody.Entry[0].Changes[0].Value.Messages[0].Text.Body
		} else if messageType == "reaction" {
			emoji = messageBody.Entry[0].Changes[0].Value.Messages[0].Reaction.Emoji
			fmt.Println("Emoji:", emoji)
		} else if messageType == "image" {
			messageBodyText = messageBody.Entry[0].Changes[0].Value.Messages[0].Image.Caption
			id = messageBody.Entry[0].Changes[0].Value.Messages[0].Image.ID
		} else if messageType == "video" {
			messageBodyText = messageBody.Entry[0].Changes[0].Value.Messages[0].Video.Caption
			id = messageBody.Entry[0].Changes[0].Value.Messages[0].Video.ID
		} else if messageType == "audio" {
			messageBodyText = messageBody.Entry[0].Changes[0].Value.Messages[0].Audio.Caption
			id = messageBody.Entry[0].Changes[0].Value.Messages[0].Audio.ID
		} else if messageType == "document" {
			messageBodyText = messageBody.Entry[0].Changes[0].Value.Messages[0].Document.Caption
			id = messageBody.Entry[0].Changes[0].Value.Messages[0].Document.ID
		}

		// Check if the message body is not empty
		if messageBodyText != "" {
			// Message body exists and is not empty
			fmt.Println("Message body:", messageBodyText)
		} else {
			// Message body is empty
			fmt.Println("Message body is empty.")
		}
	} else if len(messageBody.Entry) > 0 &&
		len(messageBody.Entry[0].Changes) > 0 &&
		messageBody.Entry[0].Changes[0].Value.Statuses != nil {
		messageStatusID = messageBody.Entry[0].Changes[0].Value.Statuses[0].ID
		messageStatus = messageBody.Entry[0].Changes[0].Value.Statuses[0].Status
		recipentId = messageBody.Entry[0].Changes[0].Value.Statuses[0].RecipientID

		fmt.Println("Status:", messageStatus)
	} else {
		// Neither Value.Messages nor Value.Statuses are present
		fmt.Println("No message body or status information available.")
	}

	switch messageType {
	case "text":
		TextMessage(ctx, from, phoneNumber, messageBodyText, profileName, msgID)
	case "reaction":
		ReactionMessage(ctx, from, phoneNumber, emoji, profileName, msgID, emoji)
	case "image":
		ImageMessage(ctx, from, phoneNumber, id, profileName, msgID, messageBodyText)
	case "video":
		VideoMessage(ctx, from, phoneNumber, id, profileName, msgID, messageBodyText)
	case "audio":
		AudioMessage(ctx, from, phoneNumber, id, profileName, msgID, messageBodyText)
	case "document":
		DocumentMessage(ctx, from, phoneNumber, id, profileName, msgID, messageBodyText)
	}

	UpdateMessageStatus(ctx, messageStatusID, messageStatus, recipentId)
	// Broadcast the message to all connected clients
	recipientIdentifier := to
	broadcastMessage(messageBody, recipientIdentifier)
}

func WebSocketHandler(ctx *gin.Context) {
	// Upgrade the connection to a WebSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Could not open websocket connection")
		return
	}

	// Assuming you can get the user identifier from the request or message
	userIdentifier := ctx.Query("user") // Example: get user from query parameter
	clients[userIdentifier] = conn
	defer func() {
		delete(clients, userIdentifier)
		conn.Close()
	}()

	for {
		// Read message from the WebSocket connection
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var message interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			continue // Ignore invalid messages
		}
		broadcastMessage(message, userIdentifier) // Pass the user identifier
	}
}

func broadcastMessage(message interface{}, recipientIdentifier string) {
	for userIdentifier, client := range clients {
		// Assuming message has a field "from" to check against the recipient
		if msgMap, ok := message.(map[string]interface{}); ok {
			if from, exists := msgMap["from"]; exists && from == recipientIdentifier {
				msg, err := json.Marshal(message)
				if err != nil {
					continue // Ignore if marshaling fails
				}
				err = client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					client.Close()
					delete(clients, userIdentifier)
				}
			}
		}
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
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)
	chatId = chat.ID

	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageType: "text",
			LastMessageBody: models.Body{
				Text: messageBody,
			},
			MessageId:   messageId,
			UserID:      replyUser.UserId,
			Status:      "Received",
			UnreadCount: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"unreadCount": chat.UnreadCount + 1, "lastMessageBody": models.Body{
			Text: messageBody,
		}, "messageId": messageId, "messageType": "text", "status": "Received", "updatedAt": time.Now()}})
	}

	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, messageBody, "", "", "text")

}

func ReactionMessage(ctx *gin.Context, from, to, messageBody, profileName, messageId, emoji string) {
	var chatId interface{}
	var replyUser models.ReplyUser
	var chat models.Chat
	var users models.User
	ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	chatId = replyUser.ID
	if replyUser.UserId == "" {
		userId := generateRandom()
		ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)
	chatId = chat.ID
	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageType: "reaction",
			LastMessageBody: models.Body{
				Text: emoji,
			},
			MessageId:   messageId,
			UserID:      replyUser.UserId,
			UnreadCount: 1,
			Status:      "Received",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"unreadCount": chat.UnreadCount + 1, "lastMessageBody": models.Body{
			Text: emoji,
		}, "messageId": messageId, "messageType": "reaction", "status": "Received", "updatedAt": time.Now()}})
	}
	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, emoji, "", "", "reaction")
	/*
		message := models.Message{
			MessageId:   messageId,
			From:        from,
			To:          to,
			MesaageType: "reaction",
			MessageBody: models.Body{
				Text: emoji,
			},
			ChatId:        chatId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			MessageStatus: "Received",
		}
		messageCollection.InsertOne(context.TODO(), message)
	*/
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
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
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
	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageType: "image",
			LastMessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "image/jpeg",
			},
			UserID:      replyUser.UserId,
			MessageId:   messageId,
			UnreadCount: 1,
			Status:      "Received",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"unreadCount": chat.UnreadCount + 1, "lastMessageBody": models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "image/jpeg",
		}, "messageId": messageId, "messageType": "image", "status": "Received", "updatedAt": time.Now()}})
	}

	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, caption, file, "", "image")
	/*
		message := models.Message{
			MessageId:   messageId,
			From:        from,
			To:          to,
			MesaageType: "image",
			MessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "image/jpeg",
			},
			ChatId:        chatId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			MessageStatus: "Received",
		}
		messageCollection.InsertOne(context.TODO(), message)*/
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
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
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
	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageId:   messageId,
			MessageType: "video",
			LastMessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "video/mp4",
			},
			UserID:      replyUser.UserId,
			UnreadCount: 1,
			Status:      "Received",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"lastMessageBody": models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "video/mp4",
		}, "messageId": messageId, "messageType": "video", "unreadCount": chat.UnreadCount + 1, "seenStatus": false, "updatedAt": time.Now()}})
	}

	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, caption, file, "", "video")
	/*
		message := models.Message{
			MessageId:   messageId,
			From:        from,
			To:          to,
			MesaageType: "video",
			MessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "video/mp4",
			},
			ChatId:        chatId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			MessageStatus: "Received",
		}
		messageCollection.InsertOne(context.TODO(), message) */

}

func AudioMessage(ctx *gin.Context, from, to, mediaId, profileName, messageId, caption string) {
	var chatId interface{}
	var replyUser models.ReplyUser
	var chat models.Chat
	var users models.User
	ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	chatId = replyUser.ID
	if replyUser.UserId == "" {
		userId := generateRandom()
		ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)
	chatId = chat.ID
	url, token, err := GetUrl(ctx, to, mediaId)
	if err != nil {
		return
	}
	file, err := DownLoadFile(ctx, url.Url, token, ".mp3")
	if err != nil {
		return
	}
	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageId:   messageId,
			MessageType: "audio",
			LastMessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "audio/mp3",
			},
			UserID:      replyUser.UserId,
			UnreadCount: 1,
			Status:      "Received",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"unreadCount": chat.UnreadCount + 1, "lastMessageBody": models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "audio/mp3",
		}, "messageId": messageId, "messageType": "audio", "status": "Received", "updatedAt": time.Now()}})
	}

	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, caption, file, "", "audio")
	/*
		message := models.Message{
			MessageId:   messageId,
			From:        from,
			To:          to,
			MesaageType: "audio",
			MessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "audio/mp3",
			},
			ChatId:        chatId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			MessageStatus: "Received",
		}
		messageCollection.InsertOne(context.TODO(), message)
	*/
}

func DocumentMessage(ctx *gin.Context, from, to, mediaId, profileName, messageId, caption string) {
	var chatId interface{}
	var replyUser models.ReplyUser
	var chat models.Chat
	var users models.User
	ReplyUserCollection.FindOne(context.TODO(), bson.M{"phoneNumber": from}).Decode(&replyUser)
	chatId = replyUser.ID
	if replyUser.UserId == "" {
		userId := generateRandom()
		ReplyUserCollection.InsertOne(context.TODO(), models.ReplyUser{PhoneNumber: from, UserId: userId, UserName: profileName})
		replyUser.UserId = userId
	}
	chatCollection.FindOne(context.TODO(), bson.M{"from": from}).Decode(&chat)
	userCollection.FindOne(context.TODO(), bson.M{"phoneNo": to}).Decode(&users)
	chatId = chat.ID
	url, token, err := GetUrl(ctx, to, mediaId)
	if err != nil {
		return
	}
	file, err := DownLoadFile(ctx, url.Url, token, ".pdf")
	if err != nil {
		return
	}
	if chat.From != from {
		user := models.Chat{
			UserName:    profileName,
			CreatedBy:   profileName,
			From:        from,
			To:          to,
			MessageType: "document",
			LastMessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "document/pdf",
			},
			MessageId:   messageId,
			UserID:      replyUser.UserId,
			UnreadCount: 1,
			Status:      "Received",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
		}
		data, _ := chatCollection.InsertOne(context.TODO(), user)
		chatId = data.InsertedID

	} else {
		chatCollection.UpdateOne(context.TODO(), bson.M{"from": from}, bson.M{"$set": bson.M{"unreadCount": chat.UnreadCount + 1, "lastMessageBody": models.Body{
			Text:     caption,
			Url:      file,
			MimeType: "document/pdf",
		}, "messageId": messageId, "messageType": "document", "status": "Received", "updatedAt": time.Now()}})
	}
	InsertIncommingMessageIntoDB(ctx, chatId, messageId, from, to, caption, file, "", "audio")
	/*
		message := models.Message{
			MessageId:   messageId,
			From:        from,
			To:          to,
			MesaageType: "document",
			MessageBody: models.Body{
				Text:     caption,
				Url:      file,
				MimeType: "document/pdf",
			},
			ChatId:        chatId,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			MessageStatus: "Received",
			IsActive:      true,
		}
		messageCollection.InsertOne(context.TODO(), message)
	*/
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
			fmt.Println(err)
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

func UpdateMessageStatus(ctx *gin.Context, messageId, status, recipentId string) {
	messageCollection.UpdateOne(context.TODO(), bson.M{"messageId": messageId}, bson.M{"$set": bson.M{"messageStatus": status}})
	chatCollection.UpdateOne(context.TODO(), bson.M{"messageId": messageId}, bson.M{"$set": bson.M{"readStatus": status, "updatedAt": time.Now()}})
}

func InsertIncommingMessageIntoDB(ctx *gin.Context, chatId interface{}, messageId, from, to, messageBody, link, parentMessageId, messageType string) (*mongo.InsertOneResult, error) {
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
		MessageId:     messageId,
		From:          from,
		To:            to,
		MesaageType:   messageType,
		MessageBody:   Body,
		ChatId:        chatId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		MessageStatus: "Received",
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
