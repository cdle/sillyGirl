package im

type Sender interface {
	GetUserID() int
	GetChatID() int
	GetImType() string
	GetMessageID() int
	GetUsername() string
	IsReply() bool
	GetReplySenderUserID() int
	GetRawMessage() interface{}
	SetMatch([]string)
	SetAllMatch([][]string)
	GetMatch() []string
	GetAllMatch() [][]string
	Get(...int) string
	GetContent() string
	IsAdmin() bool
	IsMedia() bool
	Reply(interface{}) error
}

type Config struct {
	Type         string
	Masters      []int
	Groups       []int
	Token        string
	CustomConfig string
}
