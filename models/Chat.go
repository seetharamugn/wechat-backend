package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Chat struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserName        string             `bson:"userName" json:"userName"`
	CreatedBy       string             `bson:"createdBy" json:"createdBy"`
	Status          string             `bson:"status" json:"status"`
	LastMessageBody Body               `bson:"lastMessageBody" json:"lastMessageBody"`
	MessageType     string             `bson:"messageType" json:"messageType"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
