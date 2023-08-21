package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
	"os"
	"strings"
	"time"
)

func TokenControllers(ctx *gin.Context) {
	var user models.User
	if ctx.Bind(&user) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read the input",
		})
		ctx.Abort()
		return
	}
	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email and password is required"})
		ctx.Abort()
		return
	}
	services.CreatTokenService(ctx, user)

}

func ValidateAccessToken(c *gin.Context) {
	signature := os.Getenv("JWT_SECRET_KEY")
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing authorization header",
			Data:       nil,
		})
		c.Abort()
		return
	}

	tokenString := strings.Replace(header, "Bearer ", "", 1)

	// Parse and verify the access token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signature), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, Dao.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid Token!",
			Data:       nil,
		})
		c.Abort()
		return
	}

	// Check if the token is an access token and not expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, Dao.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid Token!",
			Data:       nil,
		})
		c.Abort()
		return
	}
	expiredAt := int64(claims["expiredIn"].(float64))
	if time.Now().Unix() > expiredAt {
		c.JSON(http.StatusUnauthorized, Dao.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Expired Token!",
			Data:       nil,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Valid Token!",
		Data:       nil,
	})
	return

}
