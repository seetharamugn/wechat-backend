package models

type ReplyUser struct {
	UserId      int    `bson:"userId" json:"userId"`
	UserName    string `bson:"userName" json:"userName"`
	Profile     string `bson:"profile" json:"profile"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
}
