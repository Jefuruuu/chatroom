package messages

type IMessageRepo interface {
	Save(userId int, userName string, content string) error
	Get(messageId int) (string, error)
	List(number int) (*[]string, error)
}
