package models

import (
	"time"
)

type Message struct {
	MessageId     string      `bson:"messageId" json:"messageId"`
	From          string      `bson:"from" json:"from"`
	To            string      `bson:"to" json:"to"`
	MesaageType   string      `bson:"messageType" json:"messageType"`
	MessageBody   Body        `bson:"messageBody" json:"messageBody"`
	ChatId        interface{} `bson:"chatId,d" json:"chatId"`
	IsActive      bool        `bson:"isActive" json:"isActive"`
	ReadStatus    string      `bson:"readStatus" json:"readStatus"`
	TemplateId    string      `bson:"templateId" json:"templateId"`
	ParentId      string      `bson:"parentMessageId" json:"parentMessageId"`
	MessageStatus string      `bson:"messageStatus" json:"messageStatus"`
	CreatedAt     time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time   `bson:"updatedAt" json:"updatedAt"`
}

type Body struct {
	Text     string `json:"text"`
	Url      string `json:"url"`
	MimeType string `json:"mime_type"`
}

type MessageBody struct {
	UserId          string  `json:"userId"`
	MessageType     string  `json:"messageType"`
	MessageTo       string  `json:"messageTo"`
	MessageBody     string  `json:"messageBody"`
	MessageId       string  `json:"messageId"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	LocationAddress string  `json:"LocationAddress"`
	LocationName    string  `json:"locationName"`
	File            Body    `json:"file"`
}
