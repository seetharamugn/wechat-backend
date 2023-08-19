package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
	"strconv"
)

func CreateUser(c *gin.Context) {
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
	resp, _ := services.CreateUser(c, user)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User created successfully",
		Data:       resp,
	})

}

func GetUser(c *gin.Context) {
	userId := c.Query("userId")
	newUserId, _ := strconv.Atoi(userId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Required UserId",
			Data:       nil,
		})
		c.Abort()
		return
	}
	resp, _ := services.GetUser(c, newUserId)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User get successfully",
		Data:       resp,
	})
}

func Update(c *gin.Context) {
	userId := c.Query("userId")
	id, _ := strconv.Atoi(userId)
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.UpdateUser(c, id, body)
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

func Delete(c *gin.Context) {
	userId := c.Query("userId")
	loginUserId, _ := strconv.Atoi(userId)
	result, _ := services.DeleteUser(c, loginUserId)
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
		Data:       result,
	})
}
