package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var tokenCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "tokens")

func CreateToken(ctx *gin.Context, user models.User) (interface{}, error) {
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"username": user.Username},
	}}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid username or password",
		})
		ctx.Abort()
		return "", err
	}
	CheckPasswordHash(user.Password, existingUser.Password)
	AccessToken := GenerateAccessToken(ctx, existingUser.UserId, time.Now().Add(time.Hour*24*1).Unix())
	RefreshToken := GenerateAccessToken(ctx, existingUser.UserId, time.Now().Add(time.Hour*24*30).Unix())
	token := models.Token{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		AtExpires:    time.Now().Add(time.Hour * 24 * 1).Unix(),
		RtExpires:    time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	response := map[string]interface{}{
		"userName":    existingUser.Username,
		"firstName":   existingUser.FirstName,
		"lastName":    existingUser.LastName,
		"userId":      existingUser.UserId,
		"accessToken": token.AccessToken,
	}
	filter := bson.M{"userId": existingUser.UserId}
	update := bson.D{
		{"$set", bson.D{
			{"accessToken", token.AccessToken},
			{"refreshToken", token.RefreshToken},
			{"atExpires", time.Now().Add(time.Hour * 24 * 1).Unix()},
			{"rtExpires", time.Now().Add(time.Hour * 24 * 30).Unix()},
		}},
	}
	opts := options.Update().SetUpsert(true)
	_, err = tokenCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
	return response, err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateAccessToken(ctx *gin.Context, userId int, unix int64) string {
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredIn": unix,
	})

	AccessTokenString, err := AccessToken.SignedString([]byte("secret"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Could not generate token"})
		return ""
	}
	return AccessTokenString
}
