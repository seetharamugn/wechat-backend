package models

type Body struct {
	UserId      int    `json:"userId"`
	MessageTo   string `json:"messageTo"`
	MessageBody string `json:"messageBody"`
}

type TextMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}
type Text struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}
type ImageMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Image            Image  `json:"image"`
}
type Image struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

type LocationMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Location         Location `json:"location"`
}
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

type TextReply struct {
	MessagingProduct string `json:"messaging_product"`
	Context          struct {
		MessageId string `json:"message_id"`
	}
	To   string `json:"to"`
	Type string `json:"type"`
	Text struct {
		PreviewUrl bool   `json:"preview_url"`
		Body       string `json:"body"`
	}
}
