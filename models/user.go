package models

import "time"

type User struct {
	UserId     int       `bson:"userId" json:"userId"`
	FirstName  string    `bson:"firstName" json:"firstName"`
	LastName   string    `bson:"lastName" json:"lastName"`
	Username   string    `bson:"userName" json:"userName"`
	Password   string    `bson:"password" json:"password"`
	Email      string    `bson:"email" json:"email"`
	PhoneNo    string    `bson:"phoneNo" json:"phoneNo"`
	UserActive bool      `bson:"userActive" json:"userActive"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  time.Time `bson:"deletedAt" json:"deletedAt"`
}
