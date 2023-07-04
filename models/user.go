package models

import "time"

type User struct {
	UserId     int       `bson:"userId"`
	FirstName  string    `bson:"firstName"`
	LastName   string    `bson:"lastName"`
	Username   string    `bson:"username"`
	Password   string    `bson:"password"`
	Email      string    `bson:"email"`
	PhoneNo    string    `bson:"phoneNo"`
	UserActive bool      `bson:"userActive"`
	CreatedAt  time.Time `bson:"createdAt"`
	UpdatedAt  time.Time `bson:"updatedAt"`
	DeletedAt  time.Time `bson:"deletedAt"`
}
