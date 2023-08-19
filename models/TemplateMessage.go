package models

import "time"

type MessageTemplate struct {
	TemplateId int       `bson:"templateId" json:"templateId"`
	Name       string    `bson:"name" json:"name"`
	Content    string    `bson:"content" json:"content"`
	Footer     string    `bson:"footer" json:"footer"`
	Category   string    `bson:"category" json:"category"`
	Language   string    `bson:"language" json:"language"`
	Status     string    `bson:"status" json:"status"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time `bson:"deleted_at" json:"deleted_at"`
}
