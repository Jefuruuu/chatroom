package messages

import (
	"time"
)

type InMemoryMessageRepo struct {
	MaxKey int
	Messages map[int]MessageInfo
}

type MessageInfo struct {
	MessageId int
	UserId int
	Content string
	TimeSent time.Time
}

func NewMessageRepo() InMemoryMessageRepo {
	return InMemoryMessageRepo{
		MaxKey: -1,
		Messages: map[int]MessageInfo{},
	}
}

func (messageRepo *InMemoryMessageRepo)CheckIfMessageExist(messageId int) bool {
	_, ok := messageRepo.Messages[messageId]
	return ok
}

func (messageRepo *InMemoryMessageRepo) Save(userId int, userName string, content string) error {
	messageRepo.Messages[messageRepo.MaxKey] = MessageInfo{
		MessageId: messageRepo.MaxKey,
		UserId: userId,
		Content: content,
		TimeSent: time.Now(),
	}
	messageRepo.MaxKey = messageRepo.MaxKey + 1
	return nil
}

func (messageRepo *InMemoryMessageRepo) Get(messageId int) (string, error) {
	messageRepo.MaxKey = messageRepo.MaxKey - 1
	return messageRepo.Messages[messageId].Content ,nil
}
func (messageRepo *InMemoryMessageRepo) List(number int) (*[]string, error) {
	if messageRepo.MaxKey < number {
		number = messageRepo.MaxKey
	}
	var arr = []string{}
	for i := messageRepo.MaxKey - number; i < messageRepo.MaxKey; i = i + 1 {
		if messageGet, err := messageRepo.Get(i); err != nil {
			return nil, err
		} else {
			arr = append(arr, messageGet)
		}
	}
	return &arr, nil
}