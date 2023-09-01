package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BulkMessage struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId         string             `bson:"userId" json:"userId"`
	TemplateName   string             `bson:"templateName" json:"templateName"`
	ContactNumbers string             `bson:"contactNumbers" json:"contactNumbers"`
	ScheduledAt    time.Time          `bson:"scheduledTime" json:"scheduledTime"`
	CreateAt       time.Time          `bson:"createAt" json:"createAt"`
}
