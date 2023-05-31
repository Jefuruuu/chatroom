package tokens

type ITokenRepo interface {
	Save(userName string, token *string) error
	Get(userName string) (string, error) // return err if nothing found
	Remove(userName string) error
}
