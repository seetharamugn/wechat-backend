package Dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId      string             `json:"userId"`
	Username    string             `json:"userName"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	AccessToken string             `json:"accessToken"`
	PhoneNumber string             `json:"phoneNumber"`
}

type ResponseMessage struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaId  string `json:"wa_id"`
	}
	Messages []struct {
		Id string `json:"id"`
	}
}
