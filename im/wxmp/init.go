package wxgzh

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var wxmp = core.NewBucket("wxmp")
var u2i = core.NewBucket("wxmpu2i")
var material = core.NewBucket("wxmpMaterial")

func init() {
	file_dir := core.ExecPath + "/logs/wxmp/"
	os.MkdirAll(file_dir, os.ModePerm)
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
			sender.Wait = make(chan []interface{}, 1)
			sender.uid = u2i.GetInt(msg.FromUserName)
			if sender.uid == 0 {
				sender.uid = int(time.Now().UnixNano())
				u2i.Set(msg.FromUserName, sender.uid)
			}
			core.Senders <- sender
			end := <-sender.Wait
			ss := []string{}
			url := ""
			if len(end) == 0 {
				ss = append(ss, "无法回复该消息")
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
					url = string(item.(core.ImageUrl))
				}
			}
			mediaID := ""
			if url != "" {
				filename := file_dir + fmt.Sprint(time.Now().UnixNano()) + ".jpg"
				err := func() error {
					f, err := os.Create(filename)
					if err != nil {
						return err
					}
					rsp, err := httplib.Get(url).Response()
					_, err = io.Copy(f, rsp.Body)
					if err != nil {
						f.Close()
						return err
					}
					f.Close()
					m := officialAccount.GetMaterial()
					mediaID, _, err = m.AddMaterial(message.MsgTypeImage, filename)
					if err != nil {
						return err
					}
					material.Set(mediaID, filename)
					return nil
				}()
				if err != nil {
					ss = append(ss, err.Error())
					goto TEXT
				}
				return &message.Reply{MsgType: message.MsgTypeImage, MsgData: message.NewImage(mediaID)}
			}
		TEXT:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(strings.Join(ss, "\n\n"))}
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
	Responses []interface{}
	Wait      chan []interface{}
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
