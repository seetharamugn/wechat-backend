package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserName        string             `bson:"userName" json:"userName"`
	PhoneNumber     string             `bson:"phoneNumber" json:"phoneNumber"`
	UserID          string             `bson:"userId" json:"userId"`
	CreatedBy       string             `bson:"createdBy" json:"createdBy"`
	SeenStatus      bool               `bson:"seenStatus" json:"seenStatus"`
	LastMessageBody Body               `bson:"lastMessageBody" json:"lastMessageBody"`
	MessageType     string             `bson:"messageType" json:"messageType"`
	UnreadCount     int                `bson:"unreadCount" json:"unreadCount"`
	ReadStatus      string             `bson:"readStatus" json:"readStatus"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
