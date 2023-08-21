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

func CreateUser(ctx *gin.Context, user models.User) {
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
	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// Either username or email already exists
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "already  username or email exists",
			Data:       nil,
		})
		ctx.Abort()
		return
	} else if err != mongo.ErrNoDocuments {
		// An error occurred during the query
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user",
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

func GetUser(ctx *gin.Context, userId int) {
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userId": userId}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get user",
			Data:       nil,
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User get successfully",
		Data:       user,
	})

}

func UpdateUser(ctx *gin.Context, userId int, body models.User) (*mongo.UpdateResult, error) {
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
		FirstName:  body.FirstName,
		LastName:   body.LastName,
		Username:   body.Username,
		Email:      body.Email,
		PhoneNo:    body.PhoneNo,
		UserActive: body.UserActive,
		UpdatedAt:  time.Now(),
	}

	resp, err := userCollection.ReplaceOne(context.TODO(), bson.M{"userId": userId}, updateUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update template",
			Data:       nil,
		})
		ctx.Abort()
		return nil, err
	}
	return resp, err
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
