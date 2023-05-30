package customerrors

import "errors"

var ErrLoginInfoIncorrect = errors.New("UserName or password Incorrect")
var ErrUserAlreadyExist = errors.New("User already exist")
var ErrTokenNotValid = errors.New("Unvalid token")
var ErrSaveMessage = errors.New("Failed to save message")
var ErrKeyError = errors.New("Key does not exist")
var ErrKeyDuplicate = errors.New("Key already exist")
var ErrNullUserOrPass = errors.New("Username and password can't be blank")
var ErrUsernameNotExist = errors.New("Username doesn't exist")