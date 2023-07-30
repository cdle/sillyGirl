package common

type Sender interface {
	SetPluginID(string)
	GetPluginID() string
	GetUserID() string
	GetChatID() string
	GetBotID() string
	GetImType() string
	GetMessageID() string
	RecallMessage(...interface{})
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
	GroupKick(uid string, reject_add_request bool) error
	GroupUnkick(uid string) error
	GroupBan(uid string, duration int) error
	GroupUnban(uid string) error
	AtLast()
	UAtLast()
	IsAtLast() bool
	MessagesToSend() string
	Stop()
	SetMark(interface{})
	GetMark() interface{}
	SetExpandMessageInfo(map[string]interface{})
	GetExpandMessageInfo() map[string]interface{}
	SetVar(string, interface{})
	GetVar(string) interface{}
	SetLevel(int)
	GetLevel() int
	Event() map[string]interface{}
	Action(map[string]interface{}) (interface{}, error)
}

type FakerSenderParams struct {
	Content   string
	UserID    string
	ChatID    string
	MessageID string
}
