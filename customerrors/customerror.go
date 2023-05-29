package customerrors

import "errors"

var ErrLoginInfoNotMatch = errors.New("UserName or password doesn't match")
var ErrUserNameAlreadyExist = errors.New("UserName already exist")
var ErrTokenNotValid = errors.New("Unvalid token")
var ErrSaveMessage = errors.New("Failed to save message")
var ErrKeyError = errors.New("Key does not exist")
var ErrKeyDuplicate = errors.New("Key already exist")