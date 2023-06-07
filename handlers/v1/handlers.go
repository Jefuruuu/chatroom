package v1

import (
	"chatroom/constants"
	"chatroom/customerrors"
	"chatroom/services/storage/messages"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"chatroom/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := strings.Split(ctx.Request.Header.Get("Authorization"), " ")[1]
		fmt.Println("authHeader:", authHeader)

		jwtClaims, err := utils.ParseToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Msg": "Authenticate Failed",
			})
			return
		}
		ctx.Set("username", jwtClaims["name"])
		ctx.Next()
		}
}

func RegisterHandler(userLoginR users.IUserLoginRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerInfo utils.UserInfo
		if err := ctx.Bind(&registerInfo); err != nil {
			fmt.Println("Register Binding")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
		}

		// Check register infromation
		if err := utils.CheckRegisterInfo(registerInfo, userLoginR); err != nil {
			fmt.Println("Register CheckUserInfo")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Bcrypt password
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), constants.GenerateHashCost)

		// Save login info into userLoginR
		if err := userLoginR.Save(registerInfo.Username, passwordHash); err != nil {
			fmt.Println("Register SaveUserInfo")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Msg": fmt.Sprintf("User %s register successfully!", registerInfo.Username),
		})
	}
}

func LoginHandler(userLoginR users.IUserLoginRepo, tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context){
		var loginInfo utils.UserInfo

		// Binding username & password
		if err := ctx.Bind(&loginInfo); err != nil {
			fmt.Println("Login Binding")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Validate password
		if err := utils.ValidPassword(loginInfo, userLoginR); err != nil {
			println("Login Validate Password")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Generate and save token
		token, err := utils.NewTokenToUser(loginInfo, tokenR)
		if err != nil {
			println("Login New Token")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Msg": fmt.Sprintf("User %s login successfully!", loginInfo.Username),
			"Token": *token,
		})
	}
}

func LogoutHandler(tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, _ := ctx.Get("username")
		fmt.Println("Logout username:", username)
		// delete token in tokenR
		if err := tokenR.Remove(username.(string)); err != nil {
			println("Logout RemoveToken")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Msg": "Logout successfully",
		})
	}
}

func SendMessage(tokenR tokens.ITokenRepo, messageR messages.IMessageRepo, userId int, userName string, content string) error {
	// Validate token
	if !utils.ValidateToken(tokenR, userName) {
		return customerrors.ErrTokenNotValid
	}

	// Save message into messageR
	if err := messageR.Save(userId, userName, content); err != nil {
		return customerrors.ErrSaveMessage
	}
	return nil
}

func ListHistory(tokenR tokens.ITokenRepo, messageR messages.IMessageRepo, userName string, number int) (*[]string, error) {
	// Validate token
	if !utils.ValidateToken(tokenR, userName) {
		return nil, customerrors.ErrTokenNotValid
	}

	// Get messages from messageR
	history, err := messageR.List(number)
	if err != nil {
		fmt.Println("Failed to get chatting history")
		return nil, err
	}
	return history, nil
}