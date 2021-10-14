package core

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Sender interface {
	GetUserID() interface{}
	GetChatID() interface{}
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
	Reply(...interface{}) (int, error)
	Delete() error
	Disappear(lifetime ...time.Duration)
	Finish()
	Continue()
	IsContinue() bool
	Await(Sender, func(Sender) interface{}, ...interface{})
}

type Edit int
type Replace int
type Notify int
type Article []string

var E Edit
var R Replace
var N Notify

type ImageUrl string
type ImagePath string

type Faker struct {
	Message string
	Type    string
	UserID  interface{}
	BaseSender
}

func (sender *Faker) GetContent() string {
	return sender.Message
}

func (sender *Faker) GetUserID() interface{} {
	return sender.UserID
}

func (sender *Faker) GetChatID() interface{} {
	return 0
}

func (sender *Faker) GetImType() string {
	if sender.Type == "" {
		return "fake"
	}
	return sender.Type
}

func (sender *Faker) GetMessageID() int {
	return 0
}

func (sender *Faker) GetUsername() string {
	return ""
}

func (sender *Faker) IsReply() bool {
	return false
}

func (sender *Faker) GetReplySenderUserID() int {
	return 0
}

func (sender *Faker) GetRawMessage() interface{} {
	return sender.Message
}

func (sender *Faker) IsAdmin() bool {
	return true
}

func (sender *Faker) IsMedia() bool {
	return false
}

func (sender *Faker) Reply(msgs ...interface{}) (int, error) {
	rt := ""
	var n *Notify
	for _, msg := range msgs {
		switch msg.(type) {
		case []byte:
			rt = (string(msg.([]byte)))
		case string:
			rt = (msg.(string))
		case Notify:
			v := msg.(Notify)
			n = &v
		}
	}
	if rt != "" && n != nil {
		NotifyMasters(rt)
	}
	return 0, nil
}

func (sender *Faker) Delete() error {
	return nil
}

func (sender *Faker) Disappear(lifetime ...time.Duration) {

}

func (sender *Faker) Finish() {

}

type BaseSender struct {
	matches [][]string
	goon    bool
	child   Sender
}

func (sender *BaseSender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}
func (sender *BaseSender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *BaseSender) GetMatch() []string {
	return sender.matches[0]
}

func (sender *BaseSender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *BaseSender) Continue() {
	sender.goon = true
}

func (sender *BaseSender) IsContinue() bool {
	return sender.goon
}

func (sender *BaseSender) Get(index ...int) string {
	i := 0
	if len(index) != 0 {
		i = index[0]
	}
	if len(sender.matches) == 0 {
		return ""
	}
	if len(sender.matches[0]) < i+1 {
		return ""
	}
	return sender.matches[0][i]
}

func (sender *BaseSender) Delete() error {
	return nil
}

func (sender *BaseSender) Disappear(lifetime ...time.Duration) {

}

func (sender *BaseSender) Finish() {

}

func (sender *BaseSender) IsMedia() bool {
	return false
}

func (sender *BaseSender) GetRawMessage() interface{} {
	return nil
}

func (sender *BaseSender) IsReply() bool {
	return false
}

func (sender *BaseSender) GetMessageID() int {
	return 0
}

func (sender *BaseSender) GetUserID() interface{} {
	return nil
}
func (sender *BaseSender) GetChatID() interface{} {
	return nil
}
func (sender *BaseSender) GetImType() string {
	return ""
}

var TimeOutError = errors.New("指令超时")
var InterruptError = errors.New("被其他指令中断")

var waits sync.Map

type Carry struct {
	Chan    chan interface{}
	Pattern string
	Result  chan interface{}
	Sender  Sender
}

func (_ *BaseSender) Await(sender Sender, callback func(Sender) interface{}, params ...interface{}) {
	c := &Carry{}
	timeout := time.Second * 20
	var handleErr func(error)
	for _, param := range params {
		switch param.(type) {
		case string:
			c.Pattern = param.(string)
		case time.Duration:
			du := param.(time.Duration)
			if du != 0 {
				timeout = du
			}
		case func() string:
			callback = param.(func(Sender) interface{})

		case func(error):
			handleErr = param.(func(error))
		}
	}
	if callback == nil {
		return
	}
	if c.Pattern == "" {
		c.Pattern = `[\s\S]*`
	}
	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)
	key := fmt.Sprintf("u=%v&c=%v&i=%v", sender.GetUserID(), sender.GetChatID(), sender.GetImType())
	if oc, ok := waits.LoadOrStore(key, c); ok {
		oc.(*Carry).Chan <- InterruptError
	}
	select {
	case result := <-c.Chan:
		switch result.(type) {
		case Sender:
			waits.Delete(key)
			c.Result <- callback(result.(Sender)) //, nil
			return
		case error:
			waits.Delete(key)
			if handleErr != nil {
				handleErr(result.(error))
			}
			c.Result <- nil //
		}
	case <-time.After(timeout):
		waits.Delete(key)
		if handleErr != nil {
			handleErr(TimeOutError)
		}
		c.Result <- nil
	}
}
