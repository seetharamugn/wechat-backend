package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReplyUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId      string             `bson:"userId" json:"userId"`
	UserName    string             `bson:"userName" json:"userName"`
	Profile     string             `bson:"profile" json:"profile"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
}
