package wxgzh

import (
	"encoding/json"
	"fmt"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	"github.com/rixingyike/wechat"
)

var wxsv = core.NewBucket("wxsv")

func init() {
	wechat.Debug = true

	cfg := &wechat.WxConfig{
		Token:          wxsv.Get("token"),
		AppId:          wxsv.Get("app_id"),
		Secret:         wxsv.Get("app_secret"),
		EncodingAESKey: wxsv.Get("encoding_aes_key"),
	}

	app := wechat.New(cfg)
	// app.SendText("@all", "Hello,World!")

	core.Server.Any("/wxsv", func(c *gin.Context) {
		ctx := app.VerifyURL(c.Writer, c.Request)
		data, _ := json.Marshal(ctx.Msg)
		fmt.Println(string(data))
		ctx.NewText(string(data)).Reply()
	})

	// http.HandleFunc("/wxsv",
}

// type Sender struct {
// 	Message   string
// 	Responses []interface{}
// 	Wait      chan []interface{}
// 	uid       string
// 	core.BaseSender
// }

// func (sender *Sender) GetContent() string {
// 	if sender.Content != "" {
// 		return sender.Content
// 	}

// 	return sender.Message
// }

// func (sender *Sender) GetUserID() string {
// 	return sender.uid
// }

// func (sender *Sender) GetChatID() int {
// 	return 0
// }

// func (sender *Sender) GetImType() string {
// 	return "wxmp"
// }

// func (sender *Sender) GetMessageID() int {
// 	return 0
// }

// func (sender *Sender) GetUsername() string {
// 	return ""
// }

// func (sender *Sender) IsReply() bool {
// 	return false
// }

// func (sender *Sender) GetReplySenderUserID() int {
// 	return 0
// }

// func (sender *Sender) GetRawMessage() interface{} {
// 	return sender.Message
// }

// func (sender *Sender) IsAdmin() bool {
// 	return strings.Contains(wxmp.Get("masters"), fmt.Sprint(sender.uid))
// }

// func (sender *Sender) IsMedia() bool {
// 	return false
// }

// func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
// 	sender.Responses = append(sender.Responses, msgs...)
// 	return 0, nil
// }

// func (sender *Sender) Delete() error {
// 	return nil
// }

// func (sender *Sender) Disappear(lifetime ...time.Duration) {

// }

// func (sender *Sender) Finish() {
// 	if sender.Responses == nil {
// 		sender.Responses = []interface{}{}
// 	}
// 	sender.Wait <- sender.Responses
// }

// func (sender *Sender) Copy() core.Sender {
// 	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
// 	return &new
// }
