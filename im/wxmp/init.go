package wxgzh

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	server "github.com/rixingyike/wechat"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var wxmp = core.NewBucket("wxmp")
var material = core.NewBucket("wxmpMaterial")
var image2Md = core.NewBucket("image2Md")
var file_dir = "logs/wxmp/"

func init() {
	os.MkdirAll(file_dir, os.ModePerm)
	if !wxmp.GetBool("isKe?", false) {
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
				if msg.Event == "subscribe" {
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(wxmp.Get("subscribe_reply", "感谢关注！"))}
				}
				sender := &Sender{}
				sender.Message = msg.Content
				sender.Wait = make(chan []interface{}, 1)
				sender.uid = fmt.Sprint(msg.FromUserName)
				sender.tp = "wxmp"
				core.Senders <- sender
				end := <-sender.Wait
				ss := []string{}
				url := ""
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
						url = string(item.(core.ImageUrl))
					}
				}
				mediaID := ""
				if url != "" && len(ss) == 0 {
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
		return
	} else {
		cfg := &server.WxConfig{
			AppId:          wxmp.Get("app_id"),
			Secret:         wxmp.Get("app_secret"),
			Token:          wxmp.Get("token"),
			EncodingAESKey: wxmp.Get("encoding_aes_key"),
			DateFormat:     "XML",
		}
		app := server.New(cfg)
		core.Pushs["wxmp"] = func(i1 interface{}, s1 string, _ interface{}, _ string) {
			app.SendText(fmt.Sprint(i1), s1)
		}
		core.Server.Any("/wx/", func(c *gin.Context) {
			ctx := app.VerifyURL(c.Writer, c.Request)
			if ctx.Msg.Event == "subscribe" {
				ctx.NewText(wxmp.Get("subscribe_reply", "感谢关注！")).Reply()
				return
			}
			sender := &Sender{}
			sender.tp = "wxmp"
			sender.Message = ctx.Msg.Content
			sender.uid = ctx.Msg.FromUserName
			sender.ctx = ctx
			core.Senders <- sender
		})
	}
}

type Sender struct {
	ctx       *server.Context
	Message   string
	Responses []interface{}
	Wait      chan []interface{}
	uid       string
	tp        string
	admin     bool
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
	return sender.tp
}

func (sender *Sender) GetMessageID() string {
	return ""
}

func (sender *Sender) GetUsername() string {
	return ""
}

func (sender *Sender) GetChatname() string {
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
	if sender.admin {
		return true
	}
	return strings.Contains(core.Bucket(sender.tp).Get("masters"), fmt.Sprint(sender.uid))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) ([]string, error) {
	if sender.ctx != nil {
		rt := ""
		for _, item := range msgs {
			switch item.(type) {
			case error:
				rt = item.(error).Error()
			case string:
				rt = item.(string)
				if sender.ctx != nil && sender.Atlast && !sender.IsFinished {
					sender.ToSendMessages = append(sender.ToSendMessages, rt)
					return []string{}, nil
				}
			case []byte:
				rt = string(item.([]byte))
			case core.ImageUrl:

			}
		}
		{
			// paths := []string{}
			rt = regexp.MustCompile(`file=[^\[\]]*,url`).ReplaceAllString(rt, "file")
			for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(rt, -1) {
				// paths = append(paths, v[1])
				if strings.HasPrefix(v[1], "http") {
					if mid := image2Md.Get(v[1]); mid != "" {
						app.SendImage(sender.GetUserID(), mid)
						continue
					}
					filename := file_dir + fmt.Sprint(time.Now().UnixNano()) + ".jpg"
					err := func() error {
						f, err := os.Create(filename)
						if err != nil {
							return err
						}
						rsp, err := httplib.Get(v[1]).Response()
						_, err = io.Copy(f, rsp.Body)
						if err != nil {
							f.Close()
							return err
						}
						f.Close()
						// m := officialAccount.GetMaterial()
						// mediaID, _, err = m.AddMaterial(message.MsgTypeImage, filename)
						md, err := app.MediaUpload("image", filename)
						if err != nil {
							return err
						}
						app.SendImage(sender.GetUserID(), md.MediaID)
						image2Md.Set(v[1], md.MediaID)
						return nil
					}()
					if err == nil {

					}
				}
				//
			}
		}
		rt = regexp.MustCompile(`\[CQ:([^\[\]]+)\]`).ReplaceAllString(rt, "")
		rt = regexp.MustCompile(`^\s`).ReplaceAllString(rt, "")
		sender.ctx.NewText(rt).Send()
		return []string{}, nil
	} else {
		sender.Responses = append(sender.Responses, msgs...)
	}
	return []string{}, nil
}

func (sender *Sender) Delete() error {
	return nil
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {
	sender.IsFinished = true
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
