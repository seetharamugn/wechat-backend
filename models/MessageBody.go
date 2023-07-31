package models

type Body struct {
	UserId      int    `json:"userId"`
	MessageTo   string `json:"messageTo"`
	MessageBody string `json:"messageBody"`
	MessageId   string `json:"messageId"`
}

// TemplateMessage Body
type TemplateMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}
type Template struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}
type Language struct {
	Code string `json:"code"`
}

// TextMessage Body
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

// ImageMessage Body
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

// VideoMessage Body
type VideoMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Video            Video  `json:"video"`
}
type Video struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

// LocationMessage Body
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

// TextReply Body
type TextReply struct {
	MessagingProduct string    `json:"messaging_product"`
	RecipientType    string    `json:"recipient_type"`
	To               string    `json:"to"`
	Context          Context   `json:"context"`
	Type             string    `json:"type"`
	Text             ReplyText `json:"text"`
}
type Context struct {
	MessageId string `json:"message_id"`
}
type ReplyText struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}
