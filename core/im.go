package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Sender interface {
	GetUserID() string
	GetChatID() int
	GetImType() string
	GetMessageID() int
	GetUsername() string
	GetChatname() string
	IsReply() bool
	GetReplySenderUserID() int
	GetRawMessage() interface{}
	SetMatch([]string)
	SetAllMatch([][]string)
	GetMatch() []string
	GetAllMatch() [][]string
	Get(...int) string
	GetContent() string
	SetContent(string)
	IsAdmin() bool
	IsMedia() bool
	Reply(...interface{}) (int, error)
	Delete() error
	Disappear(lifetime ...time.Duration)
	Finish()
	Continue()
	IsContinue() bool
	ClearContinue()
	Await(Sender, func(Sender) interface{}, ...interface{}) interface{}
	Copy() Sender
}

type Edit int
type Replace int
type Notify int
type Article []string

var E Edit
var R Replace
var N Notify

type ImageData []byte
type ImageBase64 string
type ImageUrl string
type ImagePath string
type VideoUrl string

type Faker struct {
	Message string
	Type    string
	UserID  string
	ChatID  int
	BaseSender
}

func (sender *Faker) GetContent() string {
	return sender.Message
}

func (sender *Faker) GetUserID() string {
	return sender.UserID
}

func (sender *Faker) GetChatID() int {
	return sender.ChatID
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

func (sender *Faker) GetChatname() string {
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

func (sender *Faker) Copy() Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Faker)
	return &new
}

type BaseSender struct {
	matches [][]string
	goon    bool
	child   Sender
	Content string
}

func (sender *BaseSender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}
func (sender *BaseSender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *BaseSender) SetContent(content string) {
	sender.Content = content
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

func (sender *BaseSender) ClearContinue() {
	sender.goon = false
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

func (sender *BaseSender) GetUserID() string {
	return ""
}
func (sender *BaseSender) GetChatID() int {
	return 0
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

type forGroup string

type again string

var Again again = ""

var GoAgain = func(str string) again {
	return again(str)
}

type YesOrNo string

var YesNo YesOrNo = "yeson"
var Yes YesOrNo = "yes"
var No YesOrNo = "no"

type Range []int

type Switch []string

var ForGroup forGroup

func (_ *BaseSender) Await(sender Sender, callback func(Sender) interface{}, params ...interface{}) interface{} {
	c := &Carry{}
	timeout := time.Second * 86400000
	var handleErr func(error)
	var fg *forGroup
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
		case forGroup:
			a := param.(forGroup)
			fg = &a
		}
	}
	// if callback == nil {
	// 	return nil
	// }
	if c.Pattern == "" {
		c.Pattern = `[\s\S]*`
	}
	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)

	key := fmt.Sprintf("u=%v&c=%v&i=%v", sender.GetUserID(), sender.GetChatID(), sender.GetImType())
	if fg != nil {
		key += fmt.Sprintf("&t=%v&f=true", time.Now().Unix())
	}
	if oc, ok := waits.LoadOrStore(key, c); ok {
		oc.(*Carry).Chan <- InterruptError
	}
	defer func() {
		waits.Delete(key)
	}()
	for {
		select {
		case result := <-c.Chan:
			switch result.(type) {
			case Sender:
				s := result.(Sender)
				if callback == nil {
					return s.GetContent()
				}
				result := callback(s)
				if v, ok := result.(again); ok {
					if v == "" {
						c.Result <- nil
					} else {
						c.Result <- string(v)
					}
				} else if _, ok := result.(YesOrNo); ok {
					if strings.ToLower(s.GetContent()) == "y" {
						return Yes
					}

					if strings.ToLower(s.GetContent()) == "n" {
						return No
					}
					c.Result <- "Y or n ?"
				} else if vv, ok := result.(Switch); ok {
					ct := s.GetContent()
					for _, v := range vv {
						if ct == v {
							return v
						}
					}
					c.Result <- fmt.Sprintf("请从%s中选择一个。", strings.Join(vv, "、"))
				} else if vv, ok := result.(Range); ok {
					ct := s.GetContent()
					n := Int(ct)
					if fmt.Sprint(n) == ct {
						if (n >= vv[0]) && (n <= vv[1]) {

							return n
						}
					}
					c.Result <- fmt.Sprintf("请从%d~%d中选择一个整数。", vv[0], vv[1])
				} else {
					c.Result <- result
					return nil
				}
			case error:
				if handleErr != nil {
					handleErr(result.(error))
				}
				c.Result <- nil
				return nil
			}
		case <-time.After(timeout):
			if handleErr != nil {
				handleErr(TimeOutError)
			}
			c.Result <- nil
			return nil
		}
	}
}
