package main

import (
	"chatroom/handlers"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// middleware -> token validation helper
	router := gin.Default()
	
	userRepo := users.NewLoginRepo()
	tokenRepo := tokens.NewTokenRepo()

	router.POST("/register", func(ctx *gin.Context) {
		var loginInfo handlers.LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err,
			})
		}
		if err := handlers.Register(&userRepo, loginInfo.Username, loginInfo.Password); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Msg": fmt.Sprintf("User %s register successfully", loginInfo.Username), 
			})
		}
	})
	router.POST("/login", handlers.LoginHandler(&userRepo, &tokenRepo))
	
	authorized := router.Group("/auth")
	authorized.Use(handlers.Authenticate(&userRepo))
	authorized.POST("/auth/logout", handlers.LogoutHandler(&tokenRepo))
	
	router.Run()
}
