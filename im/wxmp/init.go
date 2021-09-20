package wxgzh

import (
	"fmt"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var wxmp = core.NewBucket("wxmp")
var u2i = core.NewBucket("wxmpu2i")

func init() {
	core.Server.Any("/wx/", func(c *gin.Context) {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &offConfig.Config{
			AppID:          wxmp.Get("app_id"),
			AppSecret:      wxmp.Get("app_secret"),
			Token:          wxmp.Get("token"),
			EncodingAESKey: wxmp.Get("encoding_aes_key"),
			Cache:          memory,
		}
		officialAccount := wc.GetOfficialAccount(cfg)
		server := officialAccount.GetServer(c.Request, c.Writer)
		server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
			sender := &Sender{}
			sender.Message = msg.Content
			fmt.Println(sender.Message)
			sender.Wait = make(chan string, 1)
			sender.uid = u2i.GetInt(msg.FromUserName)
			if sender.uid == 0 {
				sender.uid = int(time.Now().UnixNano())
				u2i.Set(msg.FromUserName, sender.uid)
			}
			core.Senders <- sender
			end := <-sender.Wait
			if end == "" {
				return nil
			}
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(end)}
		})
		err := server.Serve()
		if err != nil {
			return
		}
		server.Send()
	})
}

type Sender struct {
	Message   string
	matches   [][]string
	Responses []string
	Wait      chan string
	uid       int
}

func (sender *Sender) GetContent() string {
	return sender.Message
}

func (sender *Sender) GetUserID() int {
	return sender.uid
}

func (sender *Sender) GetChatID() int {
	return 0
}

func (sender *Sender) GetImType() string {
	return "wxmp"
}

func (sender *Sender) GetMessageID() int {
	return 0
}

func (sender *Sender) GetUsername() string {
	return ""
}

func (sender *Sender) IsReply() bool {
	return false
}

func (sender *Sender) GetReplySenderUserID() int {
	return 0
}

func (sender *Sender) GetRawMessage() interface{} {
	return sender.Message
}

func (sender *Sender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}
func (sender *Sender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *Sender) GetMatch() []string {
	return sender.matches[0]
}

func (sender *Sender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *Sender) Get(index ...int) string {
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

func (sender *Sender) IsAdmin() bool {
	return strings.Contains(wxmp.Get("masters"), fmt.Sprint(sender.uid))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	fmt.Println(msgs...)
	for _, item := range msgs {
		switch item.(type) {
		case string:
			sender.Responses = append(sender.Responses, item.(string))
		case []byte:
			sender.Responses = append(sender.Responses, string(item.([]byte)))
		}
	}
	return 0, nil
}

func (sender *Sender) Delete() error {
	return nil
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {
	sender.Wait <- strings.Join(sender.Responses, "\n")
}
