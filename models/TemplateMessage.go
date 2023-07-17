package models

import "time"

type MessageTemplate struct {
	TemplateId        int       `json:"TemplateId"`
	TemplateName      string    `json:"name"`
	TemplateContent   string    `json:"content"`
	TemplateFooter    string    `json:"footer"`
	TemplateCategory  string    `json:"category"`
	TemplateLanguage  string    `json:"language"`
	TemplateStatus    string    `json:"status"`
	TemplateCreatedAt time.Time `json:"created_at"`
	TemplateUpdatedAt time.Time `json:"updated_at"`
	TemplateDeletedAt time.Time `json:"deleted_at"`
}
