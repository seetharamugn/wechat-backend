package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "users")

func CreateUser(ctx *gin.Context, user models.User) (*mongo.InsertOneResult, error) {
	password, _ := HashPassword(user.Password)
	userId := GenerateRandom()
	newUser := models.User{
		UserId:    userId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Password:  password,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if username or email already exist
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"username": user.Username},
		{"email": user.Email},
	}}).Decode(&existingUser)
	if err == nil {
		// Either username or email already exists
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "already  username or email exists",
			Data:       err.Error(),
		})
		ctx.Abort()
		return nil, err
	} else if err != mongo.ErrNoDocuments {
		// An error occurred during the query
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user",
			Data:       err.Error(),
		})
		ctx.Abort()
		return nil, err
	}

	resp, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user",
			Data:       err.Error(),
		})
		ctx.Abort()
		return nil, err
	}

	return resp, nil
}

func GetUser(ctx *gin.Context, userId string) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": userId}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user",
			Data:       err.Error(),
		})
		ctx.Abort()
		return user, err
	}
	return user, nil
}

func UpdateUser(id int, body models.User) (*mongo.UpdateResult, error) {
	update := bson.D{
		{"$set", bson.D{
			{"firstName", body.FirstName},
			{"lastName", body.LastName},
			{"email", body.Email},
			{"phoneNo", body.PhoneNo},
			{"userActive", body.UserActive},
			{"updatedAt", time.Now()},
		}},
	}

	// Define the filter to identify the user to update
	filter := bson.D{{"userId", id}}

	// Perform the update operation
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	return result, err
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateRandom() int {
	randNumber := 10000000 + rand.Intn(99999999-10000000)
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": randNumber}).Decode(&user)
	if err != nil {
		return randNumber
	}
	return GenerateRandom()
}

func DeleteUser(ctx *gin.Context, userId string) (*mongo.DeleteResult, error) {
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"templateId": userId}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user",
			Data:       err.Error(),
		})
		return nil, err
	}
	resp, err := userCollection.DeleteOne(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete user",
			Data:       err.Error(),
		})
		return nil, err
	}
	return resp, nil
}
