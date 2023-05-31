package users

import "chatroom/customerrors"

type InMemoryLoginRepo struct {
	MaxKey int
	Users  map[string]UserLoginData
}

type UserLoginData struct {
	Id       int
	UserName string
	Password string
}

func NewLoginRepo() InMemoryLoginRepo {
	return InMemoryLoginRepo{
		MaxKey: -1,
		Users:  map[string]UserLoginData{},
	}
}

func (loginRepo *InMemoryLoginRepo) CheckIfKeyExist(userName string) bool {
	_, ok := loginRepo.Users[userName]
	return ok
}

func (loginRepo *InMemoryLoginRepo) Save(userName string, passwordHash []byte) error {
	loginRepo.MaxKey = loginRepo.MaxKey + 1
	if loginRepo.CheckIfKeyExist(userName) {
		return customerrors.ErrUserAlreadyExist
	}
	loginRepo.Users[userName] = UserLoginData{
		Id:       loginRepo.MaxKey,
		UserName: userName,
		Password: string(passwordHash),
	}
	return nil
}

func (loginRepo *InMemoryLoginRepo) GetPassword(userName string) (string, error) {
	if !loginRepo.CheckIfKeyExist(userName) {
		return "", customerrors.ErrUserNotExist
	}
	return loginRepo.Users[userName].Password, nil
}
