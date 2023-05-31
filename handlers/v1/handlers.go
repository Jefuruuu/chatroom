package v1

import (
	"chatroom/constants"
	"chatroom/customerrors"
	"chatroom/services/storage/messages"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"chatroom/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserInfo) Validate(userLoginR users.IUserLoginRepo) error {
	if u.Username == "" || u.Password == "" {
		return customerrors.ErrNullUserOrPass
	}

	if utils.CheckUserNameExist(userLoginR, u.Username) {
		return customerrors.ErrUserAlreadyExist
	}

	return nil
}

func (u *UserInfo) ToValidatedUser(userLoginR users.IUserLoginRepo) (*ValidRegisterUser, error) {
	u.Validate(userLoginR)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), constants.GenerateHashCost)
	if err != nil {
		return nil, customerrors.InternalServerError
	}

	return &ValidRegisterUser{
		Username:     u.Username,
		PasswordHash: passwordHash,
	}, nil
}

type ValidRegisterUser struct {
	Username     string
	PasswordHash []byte
}

// TODO: The algorithm is weird
// For example: Think about how to send a message, Can I bind two objects?
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

func RegisterHandler(ctx *gin.Context, userLoginR users.IUserLoginRepo) {
	var registerInfo UserInfo
	if err := ctx.Bind(&registerInfo); err != nil {
		fmt.Println("Register Binding")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ErrMsg": err.Error(),
		})
	}

	validatedUser, err := registerInfo.ToValidatedUser(userLoginR)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == customerrors.InternalServerError {
			statusCode = http.StatusInternalServerError
		}
		ctx.AbortWithStatusJSON(statusCode, gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}

	// Save login info into userLoginR
	// TODO: change the interface to save a validated user
	if err := userLoginR.Save(validatedUser.Username, validatedUser.PasswordHash); err != nil {
		log.Println("Register SaveUserInfo")
		// We don't expose the internal server error detail
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Msg": fmt.Sprintf("User %s register successfully!", registerInfo.Username),
	})
}

func LoginHandler(ctx *gin.Context, userLoginR users.IUserLoginRepo, tokenR tokens.ITokenRepo) {
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
	err := loginInfo.Validate(userLoginR)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}

	// Get password from storage
	passwordHash, err := userLoginR.GetPassword(loginInfo.Username)
	if err != nil {
		println("Login GetPassword")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ErrMsg": customerrors.ErrLoginInfoIncorrect,
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

func LogoutHandler(ctx *gin.Context, tokenR tokens.ITokenRepo) {
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
