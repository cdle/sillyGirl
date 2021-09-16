package qq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/go-cqhttp/coolq"
)

type Sender struct {
	Message interface{}
	matches [][]string
}

func (sender *Sender) GetContent() string {
	text := ""
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		text = coolq.ToStringMessage(sender.Message.(*message.PrivateMessage).Elements, 0, true)
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		text = coolq.ToStringMessage(m.Elements, m.GroupCode, true)
	}
	return text
}

func (sender *Sender) GetUserID() int {
	id := 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		id = int(sender.Message.(*message.PrivateMessage).Sender.Uin)
	case *message.GroupMessage:
		id = int(sender.Message.(*message.GroupMessage).Sender.Uin)
	}
	return id
}

func (sender *Sender) GetChatID() int {
	id := 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
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
	case *message.GroupMessage:
		id = int(sender.Message.(*message.GroupMessage).Id)
	}
	return id
}

func (sender *Sender) GetUsername() string {
	name := ""
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		name = sender.Message.(*message.PrivateMessage).Sender.Nickname
	case *message.GroupMessage:
		name = sender.Message.(*message.GroupMessage).Sender.Nickname
	}
	return name
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
	var sid int64 = 0
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		m := sender.Message.(*message.PrivateMessage)
		sid = m.Sender.Uin
		if m.Target == m.Sender.Uin {
			return true
		}
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		sid = m.Sender.Uin
	}
	return strings.Contains(qq.Get("masters"), fmt.Sprint(sid))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msg interface{}) error {
	switch sender.Message.(type) {
	case *message.PrivateMessage:
		m := sender.Message.(*message.PrivateMessage)
		content := ""
		switch msg.(type) {
		case string:
			content = msg.(string)
		case []byte:
			content = string(msg.([]byte))
		case *http.Response:
			data, _ := ioutil.ReadAll(msg.(*http.Response).Body)
			bot.SendPrivateMessage(m.Sender.Uin, int64(qq.GetInt("groupCode")), &message.SendingMessage{Elements: []message.IMessageElement{&coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
		}
		if content != "" {
			bot.SendPrivateMessage(m.Sender.Uin, int64(qq.GetInt("groupCode")), &message.SendingMessage{Elements: []message.IMessageElement{&message.TextElement{Content: content}}})
		}
	case *message.GroupMessage:
		m := sender.Message.(*message.GroupMessage)
		content := ""
		switch msg.(type) {
		case string:
			content = msg.(string)

		case []byte:
			content = string(msg.([]byte))
		case *http.Response:
			data, _ := ioutil.ReadAll(msg.(*http.Response).Body)
			bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&message.AtElement{Target: m.Sender.Uin}, &message.TextElement{Content: "\n"}, &coolq.LocalImageElement{Stream: bytes.NewReader(data)}}})
		}
		if content != "" {
			if strings.Contains(content, "\n") {
				content = "\n" + content
			}
			bot.SendGroupMessage(m.GroupCode, &message.SendingMessage{Elements: []message.IMessageElement{&message.AtElement{Target: m.Sender.Uin}, &message.TextElement{Content: content}}})
		}
	}
	return nil
}
