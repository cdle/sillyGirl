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
		DateFormat:     "XML",
	}
	app := wechat.New(cfg)
	core.Server.Any("/wx/", func(c *gin.Context) {
		ctx := app.VerifyURL(c.Writer, c.Request)
		// switch ctx.Msg.MsgType {
		// case wechat.TypeText:
		// 	ctx.NewText(ctx.Msg.Content).Reply() // 回复文字
		// case wechat.TypeImage:
		// 	ctx.NewImage(ctx.Msg.MediaId).Reply() // 回复图片
		// case wechat.TypeVoice:
		// 	ctx.NewVoice(ctx.Msg.MediaId).Reply() // 回复语音
		// case wechat.TypeVideo:
		// 	ctx.NewVideo(ctx.Msg.MediaId, "video title", "video description").Reply() //回复视频
		// case wechat.TypeFile:
		// 	ctx.NewFile(ctx.Msg.MediaId).Reply() // 回复文件，仅企业微信可用
		// default:
		// 	ctx.NewText("其他消息类型" + ctx.Msg.MsgType).Reply() // 回复模板消息
		// }

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
		if wxmp.GetBool("isKe?", false) {
			sender.ctx = ctx
			core.Senders <- sender
			return
		}
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
	ctx       *wechat.Context
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

	if sender.ctx != nil {
		rt := ""
		for _, item := range msgs {
			switch item.(type) {
			case error:
				rt = item.(error).Error()
			case string:
				rt = item.(string)
			case []byte:
				rt = string(item.([]byte))
			case core.ImageUrl:

			}
		}
		sender.ctx.NewText(rt).Send()
		return 0, nil
	}
	sender.Responses = append(sender.Responses, msgs...)
	return 0, nil
}

func (sender *Sender) Delete() error {
	return nil
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {
	if sender.ctx != nil {
		return
	}
	if sender.Responses == nil {
		sender.Responses = []interface{}{}
	}
	sender.Wait <- sender.Responses
}

func (sender *Sender) Copy() core.Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
	return &new
}
