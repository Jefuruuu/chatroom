package utils

import (
	"chatroom/constants"
	"chatroom/customerrors"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	name string
	iat int64
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ComparePassword(passwordInput string, passwordStored string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordStored), []byte(passwordInput)) == nil
}

func GenerateToken(userName string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userName,
		"iat": time.Now().UTC().Unix(),
	})
	tokenString, err := token.SignedString([]byte(constants.HmacSecretString))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func ValidateToken(tokenR tokens.ITokenRepo, userName string) bool {
	_, err := tokenR.Get(userName)
	return err == nil
}

func CheckUserNameExist(userLoginR users.IUserLoginRepo, userName string) bool {
	_, err := userLoginR.GetPassword(userName)
	return err == nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(interface{}, error){
		return []byte(constants.HmacSecretString), nil
	})
	if err != nil {
		fmt.Println("Parse", token, "failed:", err)
		return nil, err
	}
	return claims, nil
}

func ValidPassword(userInfo UserInfo, userLoginR users.IUserLoginRepo) (error) {
	// Check if username & password null
	if userInfo.Username == "" || userInfo.Password == "" {
		return customerrors.ErrNullUserOrPass
	}
	// Get password from storage
	passwordHash, err := userLoginR.GetPassword(userInfo.Username)
	if err != nil {
		return customerrors.ErrUserNotExist
	}
	// Check Password equal or not
	if !ComparePassword(userInfo.Password, passwordHash) {
		return customerrors.ErrLoginInfoIncorrect
	}
	return nil
}

func NewTokenToUser(userInfo UserInfo, tokenR tokens.ITokenRepo) (*string, error) {
	// Generate a token for user
	token, err := GenerateToken(userInfo.Username)
	if err != nil {
		return nil, err
	}
	// Save token to storage
	if err := tokenR.Save(userInfo.Username, token); err != nil {
		return nil, err
	}
	return token, nil
}

func CheckRegisterInfo(userInfo UserInfo, userLoginR users.IUserLoginRepo) (error) {
	// Check username and password null
	if userInfo.Username == "" || userInfo.Password == "" {
		return customerrors.ErrNullUserOrPass
	}
	// Check username exist or not
	if CheckUserNameExist(userLoginR, userInfo.Username) {
		return customerrors.ErrUserAlreadyExist
	}
	return nil
}