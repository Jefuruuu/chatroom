package main

import (
	"chatroom/handlers"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"

	"github.com/gin-gonic/gin"
)

func main() {
	// middleware -> token validation helper
	router := gin.Default()//.Group("/v1")
	v1 := router.Group("/v1")
	
	userRepo := users.NewLoginRepo()
	tokenRepo := tokens.NewTokenRepo()

	v1.POST("/register", handlers.RegisterHandler(&userRepo))
	v1.POST("/login", handlers.LoginHandler(&userRepo, &tokenRepo))

	authorized := v1.Group("/auth")
	authorized.Use(handlers.Authenticate(&tokenRepo))
	authorized.POST("/logout", handlers.LogoutHandler(&tokenRepo))
	
	router.Run()
}
