package models

type Token struct {
	Id           int
	UserId       string
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}
