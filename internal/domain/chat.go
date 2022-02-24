package domain

type Message struct {
	UserName  string `json:"user_name"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

type Messages []Message

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].CreatedAt < m[j].CreatedAt
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type Users []int64

func (us Users) Contains(id int64) bool {
	for _, item := range us {
		if item == id {
			return true
		}
	}

	return false
}

type Chat struct {
	ID    string `json:"id"`
	Users Users  `json:"users"`
}

type ChatRepo interface {
	CreateChat(id string, name string, users []int64) error
	GetHistory(id string, limit int) (Messages, error)
	SendMessage(id string, msg Message) error
	GetChat(id string) (Chat, error)
}
