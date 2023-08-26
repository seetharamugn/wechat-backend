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
	services.CreateUser(c, user)
}

func Update(c *gin.Context) {
	userId := c.Query("userId")
	var body models.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.UpdateUser(c, userId, body)
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
	services.DeleteUser(c, loginUserId)
}

func VerifyEmail(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.VerifyEmail(c, user.Email)

}

func ResetPassword(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	services.ResetPassword(c, user.Email, user.Password)

}
func GetUserDetails(c *gin.Context) {
	userId := c.Query("userId")
	services.GetUserDetails(c, userId)

}
