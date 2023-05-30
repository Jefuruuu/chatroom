package main

import (
	v1 "chatroom/handlers/v1"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"

	"github.com/gin-gonic/gin"
)

func main() {
	// middleware -> token validation helper
	router := gin.Default()//.Group("/v1")
	v1Router := router.Group("/v1")
	
	userRepo := users.NewLoginRepo()
	tokenRepo := tokens.NewTokenRepo()


	v1Router.POST("/register", v1.RegisterHandler(&userRepo))
	v1Router.POST("/login", v1.LoginHandler(&userRepo, &tokenRepo))

	authorized := v1Router.Group("/auth")
	authorized.Use(v1.Authenticate(&tokenRepo))
	authorized.POST("/logout", v1.LogoutHandler(&tokenRepo))

	router.Run()
}
