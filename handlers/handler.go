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

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(userLoginR users.IUserLoginRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginInfo LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err,
			})
		}
		if err := Register(userLoginR, loginInfo.Username, loginInfo.Password); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Username": loginInfo.Username,
				"Password": loginInfo.Password,
			})
		}
	}
}

func Register(userLoginR users.IUserLoginRepo, userName string, password string) error {
	// Check userName
	if utils.CheckUserNameExist(userLoginR, userName) {
		return customerrors.ErrUserNameAlreadyExist
	}

	// Bcrypt password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), constants.GenerateHashCost)

	// Save login info into userLoginR
	if err := userLoginR.Save(userName, passwordHash); err != nil {
		fmt.Println("Failed to save userInfo", err)
		return err
	}
	return nil
}

func LoginHandler(userLoginR users.IUserLoginRepo, tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context){
		var loginInfo LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err,
			})
		}

		if err := Login(userLoginR, tokenR, loginInfo.Username, loginInfo.Password); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Username": loginInfo.Username,
				"Password": loginInfo.Password,
			})
		}
	}
}

func Login(userLoginR users.IUserLoginRepo, tokenR tokens.ITokenRepo, userName string, password string) error {
	// Get password from userLoginRepo 
	passwordHash, err := userLoginR.GetPassword(userName)
	if err != nil {
		fmt.Println("Failed to get passeword", err)
		return err
	}
	// bcrypt.COmparepassword and hash
	// Compare password by hash
	if !utils.ComparePassword(password, passwordHash) {
		return customerrors.ErrLoginInfoNotMatch
	}
	// Generate Token for login user
	token, err := utils.GenerateToken(userName)
	if err != nil {
		fmt.Println("Failed to generate token:",err)
		return err
	}
	// Save token into tokenRepo
	if err := tokenR.Save(userName, token); err != nil {
		fmt.Println("Failed to save token", err)
		return err
	}
	return nil
}

func LogoutHandler(tokenR tokens.ITokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginInfo LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err,
			})
		}

		if err := Logout(tokenR, loginInfo.Username); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMasg": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Msg": "Logout successfully",
			})
		}

	}
}

func Logout(tokenR tokens.ITokenRepo, userName string) error {
	// validate token
	if !utils.ValidateToken(tokenR, userName) {
		return customerrors.ErrTokenNotValid
	}

	// delete token in tokenR
	if err := tokenR.Remove(userName); err != nil {
		fmt.Println("Failed to remove token", err)
		return err
	}
	return nil
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

func Authenticate(userLoginR users.IUserLoginRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(1)
		var loginInfo LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": err,
			})
		}
		fmt.Println(2)
		if !utils.CheckUserNameExist(userLoginR, loginInfo.Username) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ErrMsg": customerrors.ErrKeyError,
			})
		} else {
			ctx.Next()
		}
	}
}