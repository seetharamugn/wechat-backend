package Dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type User struct {
	Id          primitive.ObjectID `json:"id"`
	UserId      string             `json:"userId"`
	Username    string             `json:"userName"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	AccessToken string             `json:"accessToken"`
}

type UserDetails struct {
	Username    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
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
