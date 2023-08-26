package repositories

import (
	"fmt"
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
	"strconv"
	"time"
)

var userCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "users")

func CreateUser(ctx *gin.Context, user models.User) {
	password, _ := HashPassword(user.Password)
	userId := GenerateRandom()
	newUser := models.User{
		UserId:    strconv.Itoa(userId),
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
	userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)

	if existingUser.UserId != "" {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "email Already used ",
			Data:       nil,
		})
		ctx.Abort()
		return
	}

	resp, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user",
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User created successfully",
		Data:       resp,
	})
}

func UpdateUser(ctx *gin.Context, userId string, body models.User) (*mongo.UpdateResult, error) {
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": userId}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user",
			Data:       nil,
		})
		return nil, err
	}
	updateUser := models.User{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Username:    body.Username,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		UpdatedAt:   time.Now(),
	}

	filter := bson.D{{"userId", userId}}
	update := bson.D{{"$set", updateUser}}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update template",
			Data:       nil,
		})
		ctx.Abort()
		return nil, err
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

func DeleteUser(ctx *gin.Context, userId int) {
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": userId}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user",
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	resp, err := userCollection.DeleteOne(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete user",
			Data:       nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
		Data:       resp,
	})
}

func VerifyEmail(ctx *gin.Context, email string) {
	var existingUser models.User
	fmt.Println(email)
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "email verified successfully",
		Data:       nil,
	})
}

func ResetPassword(ctx *gin.Context, email, password string) {
	hashPassword, _ := HashPassword(password)
	filter := bson.D{{"email", email}}
	update := bson.D{{"$set", bson.D{{"password", hashPassword}}}}
	result, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "reset password successfully",
		Data:       result,
	})
}

func GetUserDetails(ctx *gin.Context, userId string) {
	var existingUser Dao.UserDetails
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": userId}).Decode(&existingUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "user detail fetch successfully",
		Data:       existingUser,
	})
}
