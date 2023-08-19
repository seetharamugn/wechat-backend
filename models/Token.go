package models

type Token struct {
	Id           int    `bson:"tokenId" json:"tokenId"`
	UserId       string `bson:"userId" json:"userId"`
	AccessToken  string `bson:"accessToken" json:"accessToken"`
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
	AtExpires    int64  `bson:"atExpires" json:"atExpires"`
	RtExpires    int64  `bson:"rtExpires" json:"rtExpires"`
}
