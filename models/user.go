package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId      string             `bson:"userId" json:"userId"`
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lastName" json:"lastName"`
	Username    string             `bson:"userName" json:"userName"`
	Password    string             `bson:"password" json:"password"`
	Email       string             `bson:"email" json:"email"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	UserActive  bool               `bson:"userActive" json:"userActive"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   time.Time          `bson:"deletedAt" json:"deletedAt"`
}
