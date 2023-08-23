package models

import "time"

type Chat struct {
	UserNumber  []interface{} `bson:"userNumber" json:"userNumber"`
	CreatedBy   string        `bson:"createdBy" json:"createdBy"`
	Status      string        `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
	LastMessage string        `bson:"lastMessage" json:"lastMessage"`
}
