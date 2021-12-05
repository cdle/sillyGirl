package wxgzh

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	"github.com/rixingyike/wechat"
)

var wxmp = core.NewBucket("wxmp")
var material = core.NewBucket("wxmpMaterial")

func init() {
	file_dir := "logs/wxmp/"
	os.MkdirAll(file_dir, os.ModePerm)
	wechat.Debug = true
	cfg := &wechat.WxConfig{
		AppId:          wxmp.Get("app_id"),
		Secret:         wxmp.Get("app_secret"),
		Token:          wxmp.Get("token"),
		EncodingAESKey: wxmp.Get("encoding_aes_key"),
	}
	app := wechat.New(cfg)
	core.Server.Any("/wx/", func(c *gin.Context) {
		ctx := app.VerifyURL(c.Writer, c.Request)
		// data, _ := json.Marshal(ctx.Msg)
		// fmt.Println(string(data))
		// ctx.NewText(string(data)).Reply()
		if ctx.Msg.Event == "subscribe" {
			ctx.NewText(wxmp.Get("subscribe_reply", "感谢关注！")).Reply()
			return
		}
		sender := &Sender{}
		sender.Message = ctx.Msg.Content
		sender.Wait = make(chan []interface{}, 1)
		sender.uid = ctx.Msg.FromUserName
		core.Senders <- sender
		end := <-sender.Wait
		ss := []string{}
		if len(end) == 0 {
			ss = append(ss, wxmp.Get("default_reply", "无法回复该消息"))
		}
		for _, item := range end {
			switch item.(type) {
			case error:
				ss = append(ss, item.(error).Error())
			case string:
				ss = append(ss, item.(string))
			case []byte:
				ss = append(ss, string(item.([]byte)))
			case core.ImageUrl:
				// url = string(item.(core.ImageUrl))
			}
		}
		ctx.NewText(strings.Join(ss, "\n\n")).Reply()
	})

}

type Sender struct {
	Message   string
	Responses []interface{}
	Wait      chan []interface{}
	uid       string
	core.BaseSender
}

func (sender *Sender) GetContent() string {
	if sender.Content != "" {
		return sender.Content
	}

	return sender.Message
}

func (sender *Sender) GetUserID() string {
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

func (sender *Sender) IsAdmin() bool {
	return strings.Contains(wxmp.Get("masters"), fmt.Sprint(sender.uid))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	sender.Responses = append(sender.Responses, msgs...)
	return 0, nil
}

func (sender *Sender) Delete() error {
	return nil
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {
	if sender.Responses == nil {
		sender.Responses = []interface{}{}
	}
	sender.Wait <- sender.Responses
}

func (sender *Sender) Copy() core.Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
	return &new
}
