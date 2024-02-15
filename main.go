package main

import (
	"go-backend-starter-project/api"
	"go-backend-starter-project/controllers"
	"go-backend-starter-project/initializers"
	"go-backend-starter-project/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	api.CreateRouter()

	// Basic auth
	api.Api.POST("/auth/signup", controllers.Signup)
	api.Api.POST("/auth/login", controllers.Login)
	api.Api.POST("/auth/logout", controllers.Logout)
	api.Api.GET("/auth/validate", controllers.ValidateToken)

	// Basic user
	api.Api.GET("/users", controllers.GetUsers)
	api.Api.GET("/users/:id", controllers.GetUserById)
	api.Api.DELETE("/users", middleware.RequireAuth, controllers.DeleteUser)

	// Files
	api.Router.Static("assets", "./assets")

	log.Println("Starting server...")

	err := api.Router.Run()
	if err != nil {
		log.Fatal("Router could not be created.")
	}
}
