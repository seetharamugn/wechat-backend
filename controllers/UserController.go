package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	initializers "github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to read the input",
			Data:       nil,
		})
		c.Abort()
		return
	}
	if user.FirstName == "" || user.LastName == "" || user.Username == "" || user.Password == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Required FirstName, LastName, Username, Password and Email",
		})
		c.Abort()
		return
	}
	resp, err := u.service.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User created successfully",
		Data:       resp,
	})

}

func (u *UserController) Update(c *gin.Context) {
	userId := c.Param("userId")
	id, _ := strconv.Atoi(userId)
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.service.UpdateUser(c, id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	} else {
		c.JSON(http.StatusOK, Dao.Response{
			StatusCode: http.StatusOK,
			Message:    "User updated successfully",
			Data:       result,
		})
	}
}
