package qq

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/go-cqhttp/coolq"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
)

type Sender struct {
	Message  interface{}
	matches  [][]string
	Duration *time.Duration
	deleted  bool
	core.BaseSender
}

func (sender *Sender) GetContent() string {
	if sender.Content != "" {
		return sender.Content
	}
	text := ""
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		text = coolq.ToStringMessage(sender.Message.(*message.PrivateMessage).Elements, 0, true)
	case *message.TempMessage:
		text = coolq.ToStringMessage(sender.Message.(*message.TempMessage).Elements, 0, true)
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		text = coolq.ToStringMessage(m.Elements, m.GroupCode, true)
	}
	text = strings.Replace(text, "amp;", "", -1)
	text = strings.Replace(text, "&#91;", "[", -1)
	text = strings.Replace(text, "&#93;", "]", -1)
	// sender.Reply(text)
	// text = regexp.MustCompile(`&#93;`).ReplaceAllString(text, "")
	return strings.Trim(text, " ")
}

func (sender *Sender) GetUserID() string {
	id := 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		id = int(sender.Message.(*message.PrivateMessage).Sender.Uin)
	case *message.TempMessage:
		id = int(sender.Message.(*message.TempMessage).Sender.Uin)
	case *message.GroupMessage:
		id = int(sender.Message.(*message.GroupMessage).Sender.Uin)
	}
	if id != 0 {
		return fmt.Sprint(id)
	}
	return ""
}

func (sender *Sender) GetChatID() int {
	id := 0
	switch sender.Message.(type) {
	case *message.GroupMessage:
		id = int(sender.Message.(*message.GroupMessage).GroupCode)
	}
	return id
}

func (sender *Sender) GetImType() string {
	return "qq"
}

func (sender *Sender) GetMessageID() int {
	id := 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		id = int(sender.Message.(*message.PrivateMessage).Id)
	case *message.TempMessage:
		id = int(sender.Message.(*message.TempMessage).Id)
	case *message.GroupMessage:
		id = int(sender.Message.(*message.GroupMessage).Id)
	}
	return id
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
	var sid int64 = 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		m := sender.Message.(*message.PrivateMessage)
		sid = m.Sender.Uin
		// if m.Target == m.Sender.Uin {
		// 	return true
		// }
	case *message.TempMessage:
		return false
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		sid = m.Sender.Uin
		// if bot.Client.Uin == m.Sender.Uin {
		// 	return true
		// }
	}
	id := fmt.Sprint(sid)

	if id == qq.Get("bot_id") {
		return true
	}

	for _, v := range regexp.MustCompile(`\d+`).FindAllString(qq.Get("masters"), -1) {
		if id == v {
			return true
		}
	}
	return false
}

func (sender *Sender) IsMedia() bool {
	return false
}

var dd sync.Map

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	var id int32
	if spy_on := qq.Get("spy_on"); spy_on != "" && strings.Contains(spy_on, fmt.Sprint(sender.GetChatID())) {
		return 0, nil
	}
	msg := msgs[0]
	for _, item := range msgs {
		switch item.(type) {
		case time.Duration:
			du := item.(time.Duration)
			sender.Duration = &du
		}
	}

	switch sender.Message.(type) {
	case *message.PrivateMessage:
		m := sender.Message.(*message.PrivateMessage)
		content := ""
		switch msg.(type) {
		case error:
			content = msg.(error).Error()
		case string:
			content = msg.(string)
		case []byte:
			content = string(msg.([]byte))
		case core.ImageUrl:
			data, err := httplib.Get(string(msg.(core.ImageUrl))).Bytes()
			if err != nil {
				// sender.Reply(err)
				return 0, nil
			} else {
				bot.SendPrivateMessage(m.Sender.Uin, 0, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
				//pm := // dd.Store(pm.InternalId, true)
			}
		case core.ImageData:
			bot.SendPrivateMessage(m.Sender.Uin, 0, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(msg.(core.ImageData))}}})
		case core.ImageBase64:
			data, err := base64.StdEncoding.DecodeString(string(msg.(core.ImageBase64)))
			if err != nil {
				sender.Reply(err)
				return 0, nil
			}
			bot.SendPrivateMessage(m.Sender.Uin, 0, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
		}
		if content != "" {

			bot.SendPrivateMessage(m.Sender.Uin, 0, &message.SendingMessage{Elements: bot.ConvertStringMessage(content, false)})
			//pm := // dd.Store(pm.InternalId, true)
			logs.Warn("uid=%d,gid=%d,msg=%s", m.Sender.Uin, 0)
		}

	case *message.TempMessage:
		m := sender.Message.(*message.TempMessage)
		content := ""
		switch msg.(type) {
		case error:
			content = msg.(error).Error()
		case string:
			content = msg.(string)
		case []byte:
			content = string(msg.([]byte))
		case core.ImageUrl:
			data, err := httplib.Get(string(msg.(core.ImageUrl))).Bytes()
			if err != nil {
				sender.Reply(err)
				return 0, nil
			} else {
				bot.SendPrivateMessage(m.Sender.Uin, m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
			}
		case core.ImageData:
			bot.SendPrivateMessage(m.Sender.Uin, m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(msg.(core.ImageData))}}})
		case core.ImageBase64:
			data, err := base64.StdEncoding.DecodeString(string(msg.(core.ImageBase64)))
			if err != nil {
				sender.Reply(err)
				return 0, nil
			}
			bot.SendPrivateMessage(m.Sender.Uin, m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
		}

		if content != "" {
			id = bot.SendPrivateMessage(m.Sender.Uin, m.GroupCode, &message.SendingMessage{Elements: bot.ConvertStringMessage(content, false)})
		}
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		content := ""
		switch msg.(type) {
		case error:
			content = msg.(error).Error()
		case string:
			content = msg.(string)
		case []byte:
			content = string(msg.([]byte))
		case core.ImageUrl:
			data, err := httplib.Get(string(msg.(core.ImageUrl))).Bytes()
			if err != nil {
				sender.Reply(err)
				return 0, nil
			} else {
				id = bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&message.AtElement{Target: m.Sender.Uin}, &message.TextElement{Content: " \n"}, &coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
			}
		case core.ImageData:
			id = bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&message.AtElement{Target: m.Sender.Uin}, &coolq.LocalImageElement{Stream: bytes.NewReader(msg.(core.ImageData))}}})
		case core.ImageBase64:
			data, err := base64.StdEncoding.DecodeString(string(msg.(core.ImageBase64)))
			if err != nil {
				sender.Reply(err)
				id = bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&message.AtElement{Target: m.Sender.Uin}, &coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
			}
		}
		if content != "" {
			if strings.Contains(content, "\n") {
				content = " \n" + content
			}
			id = bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: append([]message.IMessageElement{
				&message.AtElement{Target: m.Sender.Uin}}, bot.ConvertStringMessage(content, true)...)}) //
		}

	}
	if id > 0 {
		// MSG := bot.GetMessage(id)
		// dd.Store(MSG["internal-id"].(int32), true)
		// logs.Debug("send id=%d message-id=%d internal-id=%d", id, MSG["message-id"].(int32), MSG["internal-id"].(int32))
		// if m, ok := sender.Message.(*message.GroupMessage); ok {
		// 	if sender.Duration != nil {
		// 		if *sender.Duration != 0 {
		// 			go func() {
		// 				time.Sleep(*sender.Duration)
		// 				sender.Delete()

		// 				bot.Client.RecallGroupMessage(m.GroupCode, MSG["message-id"].(int32), MSG["internal-id"].(int32))
		// 			}()
		// 		} else {
		// 			sender.Delete()

		// 			bot.Client.RecallGroupMessage(m.GroupCode, MSG["message-id"].(int32), MSG["internal-id"].(int32))
		// 		}

		// 	}
		// }
	}
	return 0, nil
}

func (sender *Sender) Delete() error {
	if sender.deleted {
		return nil
	}
	switch sender.Message.(type) {
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		if err := bot.Client.RecallGroupMessage(m.GroupCode, m.Id, m.InternalId); err != nil {
			return err
		}
	}
	sender.deleted = true
	return nil
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

func (sender *Sender) Copy() core.Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
	return &new
}

func (sender *Sender) GetUsername() string {

	switch sender.Message.(type) {
	case *message.PrivateMessage:
		m := sender.Message.(*message.PrivateMessage)
		if m.Sender.Nickname == "" {
			return fmt.Sprint(m.Sender.Uin)
		}
		return m.Sender.Nickname
	case *message.TempMessage:
		m := sender.Message.(*message.TempMessage)
		if m.Sender.Nickname == "" {
			return fmt.Sprint(m.Sender.Uin)
		}
		return m.Sender.Nickname
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		if m.Sender.Nickname == "" {
			return fmt.Sprint(m.Sender.Uin)
		}
		return m.Sender.Nickname
	}
	return ""
}
