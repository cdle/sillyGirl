package tg

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Sender struct {
	Message  *tb.Message
	Duration *time.Duration
	deleted  bool
	reply    *tb.Message
	core.BaseSender
}

var tg = core.NewBucket("tg")
var b *tb.Bot
var Handler = func(message *tb.Message) {
	if message.FromGroup() {
		if groupCode := tg.GetInt("groupCode"); groupCode != 0 && groupCode != int(message.Chat.ID) {
			return
		}
	}
	core.Senders <- &Sender{
		Message: message,
	}
}

func buildClientWithProxy(addr string) (*http.Client, error) {
	if addr != "" {
		u, err := url.Parse(addr)
		if err != nil {
			panic(err)
		}
		// Patch client transport
		httpTransport := &http.Transport{Proxy: http.ProxyURL(u)}
		hc := &http.Client{Transport: httpTransport}

		return hc, nil
	}

	return nil, nil // use default
}

func init() {
	go func() {
		token := tg.Get("token")
		if token == "" {
			logs.Warn("未提供telegram机器人token")
			return
		}

		settings := tb.Settings{
			Token:  token,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
			// ParseMode: tb.ModeMarkdownV2,
                        URL: tg.Get("url"),
                        
		}
		if url := tg.Get("http_proxy"); url != "" {
			client, clientErr := buildClientWithProxy(url)
			if clientErr != nil {
				logs.Warn("telegram代理失败：%v", clientErr)
				return
			}
			settings.Client = client
		}
		var err error
		b, err = tb.NewBot(settings)

		if err != nil {
			logs.Warn("监听telegram机器人失败：%v", err)
			return
		}
		core.Pushs["tg"] = func(i interface{}, s string) {
			b.Send(&tb.User{ID: core.Int(i)}, s)
		}
		core.GroupPushs["tg"] = func(i, _ interface{}, s string) {
			paths := []string{}
			ct := &tb.Chat{ID: core.Int64(i)}
			for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(s, -1) {
				paths = append(paths, core.ExecPath+"/data/images/"+v[1])
				s = strings.Replace(s, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
			}
			s = regexp.MustCompile(`\[CQ:([^\[\]]+)\]`).ReplaceAllString(s, "")
			{
				t := []string{}
				for _, v := range strings.Split(s, "\n") {
					if v != "" {
						t = append(t, v)
					}
				}
				s = strings.Join(t, "\n")
			}
			if len(paths) > 0 {
				is := []tb.InputMedia{}
				for index, path := range paths {
					data, err := os.ReadFile(path)
					if err == nil {
						url := regexp.MustCompile("(https.*)").FindString(string(data))
						if url != "" {
							// rsp, err := httplib.Get(url).Response()
							// if err == nil {
							// 	i := &tb.Photo{File: tb.FromReader(rsp.Body)}
							// 	if index == 0 {
							// 		i.Caption = s
							// 	}
							// 	is = append(is, i)
							// }
							i := &tb.Photo{File: tb.FromURL(url)}

							if index == 0 {
								i.Caption = s
							}
							is = append(is, i)
						}
					}
				}
				b.SendAlbum(ct, is)
				return
			}
			b.Send(ct, s)
		}
		b.Handle(tb.OnPhoto, func(m *tb.Message) {
			filename := fmt.Sprint(time.Now().UnixNano()) + ".image"
			filepath := core.ExecPath + "/data/images/" + filename
			if b.Download(&m.Photo.File, filepath) == nil {
				m.Text = fmt.Sprintf(`[TG:image,file=%s]`, filename) + m.Caption
				Handler(m)
			}
		})
		// b.Handle(tb.OnSticker, func(m *tb.Message) {
		// 	buf := new(bytes.Buffer)
		// 	buf.ReadFrom(m.Sticker.FileReader)
		// 	img, err := webp.Decode(buf)
		// 	if err != nil {

		// 		return
		// 	}
		// 	buf.Reset()
		// 	imgBuf := buf
		// 	err = png.Encode(imgBuf, img)
		// 	if err != nil {
		// 		return
		// 	}
		// 	filename := fmt.Sprint(time.Now().UnixNano()) + ".image"
		// 	filepath := core.ExecPath + "/data/images/" + filename
		// 	f, err := os.Create(filepath)
		// 	if err != nil {
		// 		return
		// 	}
		// 	f.Write(imgBuf.Bytes())
		// 	f.Close()
		// 	m.Text = fmt.Sprintf(`[TG:image,file=%s]`, filename) + m.Caption
		// 	Handler(m)
		// })
		b.Handle(tb.OnText, Handler)
		logs.Info("监听telegram机器人")
		b.Start()
	}()
}

func (sender *Sender) GetContent() string {
	return sender.Message.Text
}

func (sender *Sender) GetUserID() interface{} {
	return sender.Message.Sender.ID
}

func (sender *Sender) GetChatID() interface{} {
	return int(sender.Message.Chat.ID)
}

func (sender *Sender) GetImType() string {
	return "tg"
}

func (sender *Sender) GetMessageID() int {
	return sender.Message.ID
}

func (sender *Sender) GetUsername() string {
	name := sender.Message.Sender.Username
	if name == "" {
		name = fmt.Sprint(sender.Message.Sender.ID)
	}
	return name
}

func (sender *Sender) IsReply() bool {
	return sender.Message.ReplyTo != nil
}

func (sender *Sender) GetReplySenderUserID() int {
	if !sender.IsReply() {
		return 0
	}
	return sender.Message.ReplyTo.ID
}

func (sender *Sender) GetRawMessage() interface{} {
	return sender.Message
}

func (sender *Sender) IsAdmin() bool {

	return strings.Contains(tg.Get("masters"), fmt.Sprint(sender.Message.Sender.ID))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	msg := msgs[0]
	var edit *core.Edit
	var replace *core.Replace
	for _, item := range msgs {
		switch item.(type) {
		case core.Edit:
			v := item.(core.Edit)
			edit = &v
		case core.Replace:
			v := item.(core.Replace)
			replace = &v
		case time.Duration:
			du := item.(time.Duration)
			sender.Duration = &du
		}
	}
	var rt *tb.Message
	var r tb.Recipient
	var options = []interface{}{}
	if !sender.Message.FromGroup() {
		r = sender.Message.Sender
	} else {
		r = sender.Message.Chat
		if !sender.deleted {
			options = []interface{}{&tb.SendOptions{ReplyTo: sender.Message}}
		}
	}
	var err error
	switch msg.(type) {
	case error:
		rt, err = b.Send(r, fmt.Sprintf("%v", msg), options...)
	case []byte:
		rt, err = b.Send(r, string(msg.([]byte)), options...)
	case string:
		if edit != nil && sender.reply != nil {
			if *edit == 0 {
				if sender.reply != nil {
					b.Edit(sender.reply, msg.(string))
				}
			} else {
				b.Edit(&tb.Message{
					ID: int(*edit),
				}, msg.(string))
			}
			if sender.reply == nil {
				return 0, nil
			}
			return sender.reply.ID, nil
		}
		if replace != nil {
			b.Delete(&tb.Message{
				ID: int(*edit),
			})
		}
		rt, err = b.Send(r, msg.(string), options...)
	case core.ImagePath:
		f, err := os.Open(string(msg.(core.ImagePath)))
		if err != nil {
			sender.Reply(err)
			return 0, nil
		} else {
			rts,err := b.SendAlbum(r, tb.Album{&tb.Photo{File: tb.FromReader(f)}}, options...)
			if err == nil {
				rt = &rts[0]
			}
		}
	case core.ImageUrl:
		rsp, err := httplib.Get(string(msg.(core.ImageUrl))).Response()
		if err != nil {
			sender.Reply(err)
			return 0, nil
		} else {
			rts, err := b.SendAlbum(r, tb.Album{&tb.Photo{File: tb.FromReader(rsp.Body)}}, options...)
			if err == nil {
				rt = &rts[0]
			}
		}
	}
	if err != nil {
		sender.Reply(err)
	}
	if rt != nil && sender.Duration != nil {
		if *sender.Duration != 0 {
			go func() {
				time.Sleep(*sender.Duration)
				sender.Delete()
				b.Delete(rt)
			}()
		} else {
			sender.Delete()
			b.Delete(rt)
		}
	}
	if rt != nil && sender.reply == nil {
		sender.reply = rt
	}
	if sender.reply != nil {
		return sender.reply.ID, err
	}
	return 0, nil
}

func (sender *Sender) Delete() error {
	if sender.deleted {
		return nil
	}
	msg := *sender.Message
	sender.deleted = true
	return b.Delete(&msg)
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {
	if len(lifetime) == 0 {
		sender.Duration = &core.Duration
	} else {
		sender.Duration = &lifetime[0]
	}
}

func (sender *Sender) Finish() {

}
