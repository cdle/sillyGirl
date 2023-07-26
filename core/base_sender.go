package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/utils"
)

var RegistedSenders = map[string]func() common.Sender{}

type Faker struct {
	Message string
	Type    string
	UserID  string
	ChatID  string
	Carry   chan string
	BaseSender
	Admin bool
}

//func (sender *Faker) Listen() chan string {
// 	return sender.Carry
// }

//func (sender *Faker) GetContent() string {
// 	if sender.Fsps.Content != "" {
// 		return sender.Fsps.Content
// 	}
// 	return sender.Message
// }

//func (sender *Faker) GetUserID() string {
// 	return sender.UserID
// }

//func (sender *Faker) GetBotID() string {
// 	return ""
// }

//func (sender *Faker) GetChatID() string {
// 	return sender.ChatID
// }

//func (sender *Faker) GetImType() string {
// 	if sender.Type == "" {
// 		return "fake"
// 	}
// 	return sender.Type
// }

//func (sender *Faker) GetMessageID() string {
// 	return ""
// }

//func (sender *Faker) GetUserName() string {
// 	return ""
// }

//func (sender *Faker) GetChatName() string {
// 	return ""
// }

//func (sender *Faker) IsReply() bool {
// 	return false
// }

// //func (sender *Faker) GetReplyUserID() int {
// 	return 0
// }

// //func (sender *Faker) GetRawMessage() interface{} {
// 	return sender.Message
// }

// //func (sender *Faker) IsAdmin() bool {
// 	return sender.Admin
// }

// //func (sender *Faker) IsMedia() bool {
// 	return false
// }

//func (sender *Faker) Reply(msgs ...interface{}) (string, error) {
// rt := ""
// for _, msg := range msgs {
// 	switch msg := msg.(type) {
// 	case []byte:
// 		rt = (string(msg))
// 	case string:
// 		rt = msg
// 	}
// }
// {

// 	for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(rt, -1) {
// 		// qr := qrcode2console.NewQRCode2ConsoleWithUrl(v[1], true)
// 		// defer qr.Output()
// 		rt = strings.Replace(rt, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
// 	}
// }

// if rt != "" && n != nil {
// 	NotifyMasters(rt)
// }

// if rt != "" && sender.Carry != nil {
// 	sender.Carry <- rt
// }

// 	if rt != "" && sender.Type == "terminal" {
// 		fmt.Printf("\x1b[%dm%s \x1b[0m\n", 31, rt)
// 	}
// 	return "", nil
// }

//func (sender *Faker) Delete() error {
// 	return nil
// }

//func (sender *Faker) Disappear(lifetime ...time.Duration) {

// }

//func (sender *Faker) Finish() {
// 	if sender.Carry != nil {
// 		close(sender.Carry)
// 	}
// }

//func (sender *Faker) Copy() common.Sender {
// 	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Faker)
// 	return &new
// }

//func (sender *Faker) GroupKick(uid string, reject_add_request bool) error {
// return nil
// }

//func (sender *Faker) GroupBan(uid string, duration int) error {
// 	return nil
// }

type BaseSender struct {
	matches        [][]string
	goon           bool
	Fsps           common.FakerSenderParams
	Atlast         bool
	ToSendMessages []string
	IsFinished     bool
	Duration       *time.Duration
	mark           interface{}
	params         []string
	level          int
	emf            map[string]interface{}
	id             string
	CreatedAt      time.Time
	plugin_id      string
}

func (sender *CustomSender) SetPluginID(plugin_id string) {
	sender.plugin_id = plugin_id
}
func (sender *CustomSender) GetPluginID() string {
	return sender.plugin_id
}

func (sender *CustomSender) SetLevel(l int) {
	sender.level = l
}

func (sender *CustomSender) GetLevel() int {
	return sender.level
}

func (sender *CustomSender) SetMark(mark interface{}) {
	sender.mark = mark
}

func (sender *CustomSender) GetMark() interface{} {
	return sender.mark
}

func (sender *CustomSender) SetExpandMessageInfo(emf map[string]interface{}) {
	if sender.emf == nil {
		sender.emf = map[string]interface{}{}
	}
	for key, value := range emf {
		sender.emf[key] = value
	}
}

func (sender *CustomSender) GetExpandMessageInfo() map[string]interface{} {
	if sender.emf == nil {
		sender.emf = map[string]interface{}{}
	}
	return sender.emf
}

func (sender *CustomSender) SetVar(key string, value interface{}) {
	if sender.emf == nil {
		sender.emf = map[string]interface{}{}
	}
	sender.emf[key] = value
}

func (sender *CustomSender) GetVar(key string) interface{} {
	if sender.emf == nil {
		sender.emf = map[string]interface{}{}
	}
	v, ok := sender.emf[key]
	if !ok {
		return nil
	}
	return v
}

func (sender *CustomSender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}
func (sender *CustomSender) SetParams(ss []string) {
	sender.params = ss
}
func (sender *CustomSender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *CustomSender) SetContent(content string) {
	sender.Fsps.Content = content
}

func (sender *CustomSender) SetFsps(fsps *common.FakerSenderParams) {
	sender.Fsps = *fsps
}

func (sender *CustomSender) GetMatch() []string {
	return sender.matches[0]
}

func (sender *CustomSender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *CustomSender) Continue() {
	sender.goon = true
}

func (sender *CustomSender) IsContinue() bool {
	return sender.goon
}

func (sender *CustomSender) ClearContinue() {
	sender.goon = false
}

func (sender *CustomSender) Get(i interface{}) string {
	switch i := i.(type) {
	case int:
		if len(sender.matches) == 0 {
			return ""
		}
		if len(sender.matches[0]) < i+1 {
			return ""
		}
		return sender.matches[0][i]
	case string:
		for j := range sender.params {
			if sender.params[j] == i {
				return sender.Get(j)
			}
		}
		return ""
	}
	return ""
}

func (sender *CustomSender) Delete() error {
	return nil
}

func (sender *BaseSender) Disappear(lifetime ...time.Duration) {

}

func (sender *CustomSender) Finish() {
	sender.IsFinished = true
}

func (sender *CustomSender) IsMedia() bool {
	return false
}

func (sender *CustomSender) GetRawMessage() interface{} {
	return nil
}

func (sender *CustomSender) IsReply() bool {
	return false
}

func (sender *CustomSender) Push(msg map[string]string) (string, error) {
	return "", nil
}

func (sender *CustomSender) GroupUnkick(uid string) error {
	return nil
}

func (sender *CustomSender) GetReplyUserID() int {
	return 0
}

func (sender *CustomSender) GetReplyMessageID() int {
	return 0
}

func (sender *CustomSender) AtLast() {
	sender.Atlast = true
}

func (sender *CustomSender) UAtLast() {
	sender.Atlast = false
}

func (sender *CustomSender) Stop() {
	panic("stop")
}

func (sender *CustomSender) GetID() string {
	return sender.id
}

func (sender *CustomSender) SetID() string {
	sender.id = utils.GenUUID()
	sender.CreatedAt = time.Now()
	senders.Store(sender.id, sender)
	return sender.id
}

func (sender *CustomSender) GetTime() time.Time {
	return sender.CreatedAt
}

func (sender *CustomSender) IsAtLast() bool {
	return sender.Atlast
}

func (sender *CustomSender) MessagesToSend() string {
	return strings.Join(sender.ToSendMessages, "\n")
}

var ErrorTimeOut = errors.New("指令超时")
var ErrorInterrupt = errors.New("被其他指令中断")

type Carrys struct {
	list map[int64]*Carry
	sync.RWMutex
}

func (cs *Carrys) Add(key int64, c *Carry) {
	// logs.Info("add", c.Function.Rules)
	cs.Lock()
	defer cs.Unlock()
	cs.list[key] = c
}

func (cs *Carrys) Remove(Key1 int64) {
	cs.Lock()
	defer cs.Unlock()
	for key := range cs.list {
		if key == Key1 {
			// logs.Info("rem", cs.list[key].Function.Rules)
			delete(cs.list, Key1)
		}
	}
}

func (cs *Carrys) RemoveByUUID(uuid string) {
	cs.Foreach(func(key int64, c *Carry) bool {
		if c.UUID == uuid {
			cs.Remove(key)
		}
		return true
	})
}

func (cs *Carrys) Foreach(f func(key int64, c *Carry) bool) {
	cs.RLock()
	defer cs.RUnlock()
	for key, c := range cs.list {
		f(key, c)
	}
}

var waits = map[int]*Carrys{
	1: {
		list: map[int64]*Carry{},
	},
	2: {
		list: map[int64]*Carry{},
	},
	3: {
		list: map[int64]*Carry{},
	},
	4: {
		list: map[int64]*Carry{},
	},
	5: {
		list: map[int64]*Carry{},
	},
}

type Carry struct {
	// Rules             []string
	Chan              chan interface{}
	Result            chan interface{}
	Message           common.Sender
	Function          common.Function
	RequireAdmin      bool
	AllowPlatforms    []string
	ProhibitPlatforms []string
	AllowGroups       []string
	ProhibitGroups    []string
	AllowUsers        []string
	ProhibitUsers     []string
	ListenPrivate     bool
	ListenGroup       bool
	UserID            string
	ChatID            string
	UUID              string
}

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

var listenCounter int64

func (s *CustomSender) Await(message common.Sender, callback func(common.Sender) interface{}, params ...interface{}) interface{} {
	timeout := time.Hour * 999999
	var handleErr func(error)
	var persistent = false
	var c *Carry
	for _, param := range params {
		switch param := param.(type) {
		case string:
			if param == "persistent" {
				persistent = true
			} else {
				c.Function.Rules = append(c.Function.Rules, param)
			}
		case []string:
			c.Function.Rules = append(c.Function.Rules, param...)
		case time.Duration:
			du := param
			if du != 0 {
				timeout = du
			}
		case func(error):
			handleErr = param
		case *Carry:
			c = param
		}
	}
	// c.UserID
	c.Message = message
	if len(c.Function.Rules) == 0 {
		c.Function.Rules = []string{`raw [\s\S]+`}
	}
	fmtRule(&c.Function)
	c.Chan = make(chan interface{}, 1)
	c.Result = make(chan interface{}, 1)
	key := atomic.AddInt64(&listenCounter, 1)
	// key := fmt.Sprintf("u=%v&c=%v&i=%v&t=%v&p=%v", message.GetUserID(), message.GetChatID(), message.GetImType(), atomic.LoadInt64(&listenCounter))
	// if fg != nil {
	// 	if *fg == "me" {
	// 		key += "&f=me"
	// 	} else {
	// 		key += "&f=true"
	// 	}
	// }
	waits[4-s.level].Add(key, c)
	defer func() {
		waits[4-s.level].Remove(key)
	}()
	for {
		select {
		case result := <-c.Chan:
			switch s := result.(type) {
			case common.Sender:
				if callback == nil {
					return s.GetContent()
				}
				if persistent {
					go func() {
						c.Result <- callback(s)
					}()
					continue
				}
				result := callback(s)
				if v, ok := result.(again); ok { //阻塞
					if v == "" {
						c.Result <- nil
					} else {
						c.Result <- string(v)
					}
				} else if _, ok := result.(YesOrNo); ok {
					o := strings.ToLower(regexp.MustCompile("[yYnN]").FindString(s.GetContent()))
					if o == "y" {
						return Yes
					}
					if o == "n" {
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
					c.Result <- fmt.Sprintf("请从%s中选择一个", strings.Join(vv, "、"))
				} else if vv, ok := result.(Range); ok {
					ct := s.GetContent()
					n := utils.Int(ct)
					if fmt.Sprint(n) == ct {
						if (n >= vv[0]) && (n <= vv[1]) {

							return n
						}
					}
					c.Result <- fmt.Sprintf("请从%d~%d中选择一个整数", vv[0], vv[1])
				} else {
					c.Result <- result
					return s.GetContent()
				}
			case error:
				if handleErr != nil {
					handleErr(s)
				}
				c.Result <- nil
				return nil
			}
		case <-time.After(timeout):
			if handleErr != nil {
				handleErr(ErrorTimeOut)
			}
			c.Result <- nil
			return nil
		}
	}
}
