package middleware

import (
	"fmt"
	"net/http"
	"os"

	"example.com/m/initializers"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func  RefreshTokens(c *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, _ := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])
		_, err := uuid.Parse(user.ID.String())

		if err != nil{
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		
		if int(claims["sub"].(float64)) == 1 {

			tokens, err := CreateTokens(user, c)
			if err != nil {
				return //err
			}

			//checking ip and sending email if it's needed
			if (user.IP != c.ClientIP()){
				SendEmail(user.Email)
			}

			//saving refresh token and ip
			user.RefreshToken = string(tokens["refresh_token"])
			user.IP = c.ClientIP()
			initializers.DB.Save(user)

			//send back
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Access", tokens["access_token"], 900, "", "", false, true)
			c.SetCookie("Refresh", tokens["refresh_token"], 60*15, "", "", false, true)
			c.JSON(http.StatusOK, gin.H{})
		}

		c.AbortWithStatus(http.StatusUnauthorized)
		return 
	}

	c.AbortWithStatus(http.StatusUnauthorized)
	return //err
}