package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserName        string             `bson:"userName" json:"userName"`
	From            string             `bson:"from" json:"from"`
	To              string             `bson:"to" json:"to"`
	UserID          string             `bson:"userId" json:"userId"`
	CreatedBy       string             `bson:"createdBy" json:"createdBy"`
	Status          string             `bson:"status" json:"status"`
	LastMessageBody Body               `bson:"lastMessageBody" json:"lastMessageBody"`
	MessageId       string             `bson:"messageId" json:"messageId"`
	MessageType     string             `bson:"messageType" json:"messageType"`
	UnreadCount     int                `bson:"unreadCount" json:"unreadCount"`
	ReadStatus      string             `bson:"readStatus" json:"readStatus"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsActive        bool               `bson:"isActive" json:"isActive"`
}
