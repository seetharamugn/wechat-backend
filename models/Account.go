package models

type WhatsappAccount struct {
	ID                string `json:"id"`
	UserId            string `json:"userId"`
	PhoneNumber       string `json:"phoneNumber"`
	PhoneNumberId     string `json:"phoneNumberId"`
	BusinessAccountId string `json:"businessAccpuntId"`
	Token             string `json:"token"`
	ApiVersion        string `json:"apiVersion"`
}
