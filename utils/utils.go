package utils

import (
	"chatroom/constants"
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(passwordInput string, passwordStored string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordStored), []byte(passwordInput)) == nil
}

func GenerateToken(userName string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userName,
		"iat":  time.Now().UTC().Unix(),
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
