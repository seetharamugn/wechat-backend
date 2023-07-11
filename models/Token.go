package models

type Token struct {
	Id           int    `json:"id"`
	UserId       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	AtExpires    int64  `json:"atExpires"`
	RtExpires    int64  `json:"rtExpires"`
}
