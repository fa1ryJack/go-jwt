package controllers

import (
	"net/http"

	"example.com/m/initializers"
	"example.com/m/middleware"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context){
	//get email and password from request body
	var body struct{
		Email string
		Password string
	}

	if c.Bind(&body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
		return
	}

	//create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})
}


func Signin(c *gin.Context){
	//get email and password from request body
	var body struct{
		Email string
		Password string
	}

	if c.Bind(&body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//check the user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	_, err1 := uuid.Parse(user.ID.String())

	if err1 != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No user was found. Invalid email or password?",
		})
		return
	}

	//compare the password with hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No user was found. Invalid email or password?",
		})
		return
	}

	//creating tokens
	tokens, err := middleware.CreateTokens(user, c)
	if err != nil{
		return
	}

	//saving refresh token and ip
	user.RefreshToken = string(tokens["refresh_token"])
	user.IP = c.ClientIP()
	initializers.DB.Save(user)

	//send id back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Access", tokens["access_token"], 900, "", "", false, true)
	c.SetCookie("Refresh", tokens["refresh_token"], 60*15, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context){
		user, _ := c.Get("user")


	c.JSON(http.StatusOK, gin.H{
		"validated user": user,
	})
}
