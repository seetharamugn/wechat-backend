package models

import (
	"time"
)

type Message struct {
	Id            string      `bson:"message_id" json:"message_id"`
	From          string      `bson:"from" json:"from"`
	To            string      `bson:"to" json:"to"`
	Type          string      `bson:"type" json:"type"`
	Body          string      `bson:"body" json:"body"`
	ChatId        interface{} `bson:"chat_id,d" json:"chat_id"`
	CreatedAt     time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `bson:"updated_at" json:"updated_at"`
	ReadStatus    bool        `bson:"read_status" json:"read_status"`
	TemplateId    string      `bson:"template_id" json:"template_id"`
	ParentId      string      `bson:"parent_message_id" json:"parent_message_id"`
	MessageStatus bool        `bson:"message_status" json:"message_status"`
}
