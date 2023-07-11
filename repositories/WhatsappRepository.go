package repositories

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"strconv"
)

var WaAccountCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "wa_account")

func CreateAccount(ctx *gin.Context, account models.WhatsappAccount) (string, error) {
	AccountId := GenerateAcRandom()
	newWaAccount := models.WhatsappAccount{
		AccountId:         AccountId,
		UserId:            account.UserId,
		Token:             account.Token,
		PhoneNumber:       account.PhoneNumber,
		PhoneNumberId:     account.PhoneNumberId,
		BusinessAccountId: account.BusinessAccountId,
		ApiVersion:        account.ApiVersion,
	}
	_, err := WaAccountCollection.InsertOne(context.TODO(), newWaAccount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create user",
		})
		ctx.Abort()
		return "", err
	}
	return "Registration Success", nil

}

func GetAccessToken(ctx *gin.Context, userId string) (models.WhatsappAccount, error) {
	//find the data using the userId
	userid, _ := strconv.Atoi(userId)
	var account models.WhatsappAccount
	err := WaAccountCollection.FindOne(context.TODO(), bson.M{"userId": userid}).Decode(&account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       err,
		})
		ctx.Abort()
		return models.WhatsappAccount{}, err
	}
	fmt.Println(account)
	return account, nil
}

func UpdateAccessToken(ctx *gin.Context, userId int, accessToken string) (string, error) {
	//find the data using the userId
	_, err := WaAccountCollection.UpdateOne(context.TODO(), models.WhatsappAccount{UserId: userId}, models.WhatsappAccount{Token: accessToken})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update access token",
			Data:       nil,
		})
	}
	return "Access Token Updated", nil
}
func DeleteAccessToken(ctx *gin.Context, userId int) (string, error) {
	//find the data using the userId
	_, err := WaAccountCollection.DeleteOne(context.TODO(), models.WhatsappAccount{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete access token",
			Data:       nil,
		})
	}
	return "Access Token Deleted", nil
}

func GenerateAcRandom() int {
	randNumber := 10000000 + rand.Intn(99999999-10000000)
	var waAccount models.WhatsappAccount
	err := WaAccountCollection.FindOne(context.TODO(), bson.M{"accountId": randNumber}).Decode(&waAccount)
	if err != nil {
		return randNumber
	}
	return GenerateRandom()
}
