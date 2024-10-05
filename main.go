package main

import (
	"example.com/m/controllers"
	"example.com/m/initializers"
	"example.com/m/middleware"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main(){
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/signin", controllers.Signin)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/refresh", middleware.RefreshTokens)

	r.Run()
}