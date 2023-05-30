package tokens

import (
	"chatroom/customerrors"
)

type InMemoryTokenRepo struct {
	Token map[string]string 
}

func NewTokenRepo() InMemoryTokenRepo {
	return InMemoryTokenRepo{
		Token: map[string]string{},
	}
}

func (tokenRepo *InMemoryTokenRepo)CheckIfUserExist(userName string) bool {
	_, ok := tokenRepo.Token[userName]
	return ok
}

func (tokenRepo *InMemoryTokenRepo) Save(userName string, token *string) error {
	if tokenRepo.CheckIfUserExist(userName) {
		return customerrors.ErrUserAlreadyExist
	}
	tokenRepo.Token[userName] = *token
	return nil
}

func (tokenRepo *InMemoryTokenRepo)Get(userName string) (string, error) {
	if !tokenRepo.CheckIfUserExist(userName) {
		return "", customerrors.ErrUserNotExist
	}
	return tokenRepo.Token[userName], nil
}

func (tokenRepo *InMemoryTokenRepo)Remove(userName string) error {
	if !tokenRepo.CheckIfUserExist(userName) {
		return customerrors.ErrUserNotExist
	}
	delete(tokenRepo.Token, userName)
	return nil
}