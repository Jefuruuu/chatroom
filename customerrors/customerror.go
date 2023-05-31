package customerrors

import "errors"

var ErrLoginInfoIncorrect = errors.New("UserName or password Incorrect")
var ErrUserAlreadyExist = errors.New("User already exist")
var ErrTokenNotValid = errors.New("Unvalid token")
var ErrSaveMessage = errors.New("Failed to save message")
var ErrUserNotExist = errors.New("User does not exist")
var ErrUserDuplicate = errors.New("User already exist")
var ErrNullUserOrPass = errors.New("Username and password can't be blank")
var ErrUsernameNotExist = errors.New("Username doesn't exist")
var InternalServerError = errors.New("Internal server error")
