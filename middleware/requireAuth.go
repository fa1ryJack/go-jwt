package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/m/initializers"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context){
	//get the cookie from request
	tokenString, err := c.Cookie("Access")
	if err != nil{
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	//decode and validate cookie(token)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok{

		//check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64){
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//find the user with token
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])

		_, err := uuid.Parse(user.ID.String())

		if err != nil{
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//attach to request
		c.Set("user", user)

		//continue
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	
}