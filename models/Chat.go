package models

type Chat struct {
	ChatType    string   `json:"chatType"`
	Users       []string `json:"users"`
	CreatedBy   string   `json:"createdBy"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	LastMessage string   `json:"lastMessage"`
}
