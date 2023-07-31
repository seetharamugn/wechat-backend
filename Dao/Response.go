package Dao

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type User struct {
	UserId      int    `json:"userId"`
	Username    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	AccessToken string `json:"accessToken"`
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
