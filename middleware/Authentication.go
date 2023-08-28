package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/seetharamugn/wachat/Dao"
	"net/http"
	"os"
	"strings"
	"time"
)

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
	return

}
