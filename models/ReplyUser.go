package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId      string             `bson:"userId" json:"userId"`
	UserName    string             `bson:"userName" json:"userName"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
