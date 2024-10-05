package middleware

import (
	"net/http"
	"os"
	"time"

	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateTokens(user models.User, c *gin.Context)(map[string]string, error) {
	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute*15).Unix(),
		"client_ip": c.ClientIP(),
	})

	//refresh token based on user's email
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(user.Email), bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return nil, err
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token.",
		})
		return nil, err
	}
	
	

	return map[string]string{
		"access_token":  tokenString,
		"refresh_token": string(refreshToken),
	}, nil
}