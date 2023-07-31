package models

import "time"

type Message struct {
	Id         string    `json:"message_id"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Type       string    `json:"type"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ReadStatus bool      `json:"read_status"`
	TemplateId string    `json:"template_id"`
	ParentId   string    `json:"parent_message_id"`
}
