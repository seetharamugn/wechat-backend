package models

type WhatsappAccount struct {
	AccountId         int    `bson:"accountId" json:"accountId"`
	UserId            int    `bson:"userId" json:"userId"`
	PhoneNumber       string `bson:"phoneNumber" json:"phoneNumber" `
	PhoneNumberId     int    `bson:"phoneNumberId"  json:"phoneNumberId"`
	BusinessAccountId int    `bson:"businessAccountId"  json:"businessAccountId"`
	Token             string `bson:"token" json:"token" `
	ApiVersion        string `bson:"apiVersion"  json:"apiVersion"`
}
