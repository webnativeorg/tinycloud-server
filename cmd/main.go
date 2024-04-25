package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"github.com/webnativeorg/tinycloud-server/cmd/handlers"
	"github.com/webnativeorg/tinycloud-server/cmd/middlewares"
	"github.com/webnativeorg/tinycloud-server/cmd/services"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// User routes
	r.POST("/register", services.RegisterUser)
	r.POST("/login", services.Login)
	// JWT validation
	authorized := r.Group("/")
	authorized.Use(middlewares.ValidateJWT())
	{
		authorized.GET("/users", services.GetUsers)

		// File routes
		authorized.POST("/files/upload", handlers.UploadFilesHandler)
	}

	r.Run(":" + environment.PORT)
	fmt.Println("Server running on port: ", environment.PORT)
}
