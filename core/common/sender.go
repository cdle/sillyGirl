package common

type Sender interface {
	GetUserID() string
	GetChatID() string
	GetBotID() string
	GetImType() string
	GetMessageID() string
	RecallMessage(...interface{}) error
	GetUserName() string
	GetChatName() string
	IsReply() bool
	GetReplyUserID() int
	GetReplyMessageID() int
	GetRawMessage() interface{}
	SetMatch([]string)
	SetParams([]string)
	SetAllMatch([][]string)
	GetMatch() []string
	GetAllMatch() [][]string
	Get(interface{}) string
	GetContent() string
	SetContent(string)
	SetFsps(fsps *FakerSenderParams)
	IsAdmin() bool
	IsMedia() bool
	Reply(...interface{}) (string, error)
	Push(map[string]string) (string, error)
	Delete() error
	Finish()
	Continue()
	IsContinue() bool
	ClearContinue()
	Await(Sender, func(Sender) interface{}, ...interface{}) interface{}
	Copy() Sender
	GroupKick(uid string, reject_add_request bool)
	GroupUnkick(uid string)
	GroupBan(uid string, duration int)
	GroupUnban(uid string)
	AtLast()
	UAtLast()
	IsAtLast() bool
	MessagesToSend() string
	Stop()
	SetMark(interface{})
	GetMark() interface{}
	SetLevel(int)
	GetLevel() int
}

type FakerSenderParams struct {
	Content string
	UserID  string
	ChatID  string
}
