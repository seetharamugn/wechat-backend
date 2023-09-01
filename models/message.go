package models

import (
	"time"
)

type Message struct {
	Id            string      `bson:"messageId" json:"messageId"`
	From          string      `bson:"from" json:"from"`
	To            string      `bson:"to" json:"to"`
	Type          string      `bson:"type" json:"type"`
	Body          Body        `bson:"body" json:"body"`
	ChatId        interface{} `bson:"chatId,d" json:"chatId"`
	CreatedAt     time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time   `bson:"updatedAt" json:"updatedAt"`
	DeletedAt     time.Time   `bson:"deletedAt" json:"deletedAt"`
	ReadStatus    string      `bson:"readStatus" json:"readStatus"`
	TemplateId    string      `bson:"templateId" json:"templateId"`
	ParentId      string      `bson:"parentMessageId" json:"parentMessageId"`
	MessageStatus bool        `bson:"messageStatus" json:"messageStatus"`
}

type Body struct {
	Text     string `json:"text"`
	Url      string `json:"url"`
	MimeType string `json:"mime_type"`
}

type MessageBody struct {
	UserId          int     `json:"userId"`
	MessageTo       string  `json:"messageTo"`
	MessageBody     string  `json:"messageBody"`
	MessageId       string  `json:"messageId"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	LocationAddress string  `json:"LocationAddress"`
	LocationName    string  `json:"locationName"`
}
