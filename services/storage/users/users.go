package users

type IUserLoginRepo interface {
	Save(userName string, passwordHash []byte) error
	GetPassword(userName string) (string, error)
}
