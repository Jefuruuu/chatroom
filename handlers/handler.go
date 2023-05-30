package handlers

import (
	"chatroom/constants"
	"chatroom/customerrors"
	"chatroom/services/storage/messages"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"chatroom/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Authenticate(tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInfo UserInfo
		if err := ctx.Bind(&userInfo); err != nil {
			fmt.Println("Authenticate Binding")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}
		_, err := tokenR.Get(userInfo.Username) 
		if err != nil {
			fmt.Println("Authenticate CheckToken")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

func RegisterHandler(userLoginR users.IUserLoginRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerInfo UserInfo
		if err := ctx.Bind(&registerInfo); err != nil {
			fmt.Println("Register Binding")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
		}

		// Check username and password null
		if registerInfo.Username == "" || registerInfo.Password == "" {
			fmt.Println("Register CheckUserName")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrNullUserOrPass.Error(),
			})
			return
		}

		// Check username exist or not
		if utils.CheckUserNameExist(userLoginR, registerInfo.Username) {
			fmt.Println("Register UserNameExist")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrUserAlreadyExist.Error(),
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
		var loginInfo UserInfo

		// Binding username & password
		if err := ctx.Bind(&loginInfo); err != nil {
			fmt.Println("Login Binding")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Check if username & password null
		if loginInfo.Username == "" || loginInfo.Password == "" {
			println("Login CheckUserName")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrNullUserOrPass.Error(),
			})
			return
		}

		// Get password from storage
		passwordHash, err := userLoginR.GetPassword(loginInfo.Username)
		if err != nil {
			println("Login GetPassword")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Check Password equal or not
		if !utils.ComparePassword(loginInfo.Password, passwordHash) {
			println("Login ComparePassword")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrLoginInfoIncorrect.Error(),
			})
			return
		}

		// Generate a token for user
		token, err := utils.GenerateToken(loginInfo.Username)
		if err != nil {
			println("Login GenerateToken")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		// Save token to storage
		if err := tokenR.Save(loginInfo.Username, token); err != nil {
			println("Login SaveToken")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Msg": fmt.Sprintf("User %s login successfully!", loginInfo.Username),
		})
	}
}

func LogoutHandler(tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		println("Logout GetUserInfo")
		userInfo, _ := ctx.Get("userInfo")
		username := userInfo.(UserInfo).Username

		// validate token
		if !utils.ValidateToken(tokenR, username) {
			println("Logout TokenInvalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrTokenNotValid.Error(),
			})
			return
		}

		// delete token in tokenR
		if err := tokenR.Remove(username); err != nil {
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