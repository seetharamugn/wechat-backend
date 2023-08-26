package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"
)

var tokenCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "tokens")

func CreateToken(ctx *gin.Context, user models.User) {
	signature := os.Getenv("JWT_SECRET_KEY")
	var existingUser models.User

	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	fmt.Println(user.Email, user.Password)
	flag := CheckPasswordHash(user.Password, existingUser.Password)
	if !flag {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Data:       nil,
			Message:    "Invalid email or password",
		})
		ctx.Abort()
		return
	}
	AccessToken := GenerateAccessToken(ctx, existingUser.UserId, time.Now().Add(time.Hour*24*1).Unix(), signature)
	RefreshToken := GenerateAccessToken(ctx, existingUser.UserId, time.Now().Add(time.Hour*24*30).Unix(), signature)
	token := models.Token{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		AtExpires:    time.Now().Add(time.Hour * 24 * 1).Unix(),
		RtExpires:    time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	response := Dao.User{
		Id:          existingUser.ID,
		UserId:      existingUser.UserId,
		Username:    existingUser.Username,
		FirstName:   existingUser.FirstName,
		LastName:    existingUser.LastName,
		AccessToken: token.AccessToken,
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
	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data":       response,
	})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccessToken(ctx *gin.Context, userId string, unix int64, signature string) string {
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredIn": unix,
	})

	AccessTokenString, err := AccessToken.SignedString([]byte(signature))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
			Message:    "Could not generate token"})
		return ""
	}
	return AccessTokenString
}
