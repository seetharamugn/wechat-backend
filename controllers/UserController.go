package controllers

import (
	"github.com/gin-gonic/gin"
	initializers "github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

var userCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "users")

func (u *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read the input",
		})
		c.Abort()
		return
	}
	if user.FirstName == "" || user.LastName == "" || user.Username == "" || user.Password == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "required UserId and DeviceId",
		})
		c.Abort()
		return
	}
	resp, err := u.service.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, resp)

}

func (u *UserController) Update(c *gin.Context) {
	username := c.Param("username")
	var body models.User
	filter := bson.D{{"username", username}}
	update := bson.D{{"$set", bson.D{{"firstName", body.FirstName}, {"lastName", body.LastName}, {"email", body.Email}, {"phoneNo", body.PhoneNo}, {"userActive", body.UserActive}, {"updatedAt", body.UpdatedAt}}}}
	if _, err := userCollection.UpdateOne(context.TODO(), filter, update); err != nil {
		panic(err)
	}
}
