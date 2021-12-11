package qq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var qq = core.NewBucket("qq")

// type Params struct {
// 	UserID  interface{} `json:"user_id"`
// 	Message string      `json:"message"`
// 	GroupID int         `json:"group_id"`
// }

type CallApi struct {
	Action string                 `json:"action"`
	Echo   string                 `json:"echo"`
	Params map[string]interface{} `json:"params"`
}

type sender struct {
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserID   int    `json:"user_id"`
}

type Message struct {
	Anonymous   interface{} `json:"anonymous"`
	Font        int         `json:"font"`
	GroupID     int         `json:"group_id"`
	Message     string      `json:"message"`
	MessageID   int         `json:"message_id"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	RawMessage  string      `json:"raw_message"`
	SelfID      int         `json:"self_id"`
	Sender      sender      `json:"sender"`
	SubType     string      `json:"sub_type"`
	Time        int         `json:"time"`
	UserID      int         `json:"user_id"`
}

var conns = map[string]*websocket.Conn{}
var defaultBot = ""
var ignore = qq.Get("ignore")

func init() {
	core.OttoFuncs["qq_bots"] = func(string) string {
		ss := []string{}
		for v := range conns {
			ss = append(ss, v)
		}
		return strings.Join(ss, " ")
	}
	core.Pushs["qq"] = func(i interface{}, s string, _ interface{}, botID string) {
		if botID == "" {
			botID = defaultBot
		}
		conn, ok := conns[botID]
		if !ok {
			botID = ""
			for v := range conns {
				botID = v
				break
			}
			if botID == "" {
				return
			}
		}
		conn.WriteJSON(CallApi{
			Action: "send_private_msg",
			Params: map[string]interface{}{
				"user_id": core.Int64(i),
				"message": s,
			},
		})
	}
	core.GroupPushs["qq"] = func(i, j interface{}, s string, botID string) {
		if botID == "" {
			botID = defaultBot
		}
		conn, ok := conns[botID]
		if !ok {
			return
		}
		conn.WriteJSON(CallApi{
			Action: "send_group_msg",
			Params: map[string]interface{}{
				"group_id": core.Int(i),
				"user_id":  core.Int64(j),
				"message":  s,
			},
		})
	}
	core.Server.GET("/qq/receive", func(c *gin.Context) {
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			return
		}
		botID := c.GetHeader("X-Self-ID")
		if len(conns) == 0 {
			defaultBot = botID
		} else if qq.Get("default_bot") == botID {
			defaultBot = botID
		}
		conns[botID] = ws
		if !strings.Contains(ignore, botID) {
			ignore += "&" + botID
		}
		go func() {
			for {
				_, data, err := ws.ReadMessage()
				if err != nil {
					ws.Close()
					delete(conns, botID)
					if defaultBot == botID {
						defaultBot = ""
						for v := range conns {
							defaultBot = v
							break
						}
					}
					break
				}
				// fmt.Println(string(data))
				msg := &Message{}
				json.Unmarshal(data, msg)
				if msg.MessageType != "private" && fmt.Sprint(msg.SelfID) != defaultBot {
					continue
				}
				// fmt.Println(msg)
				if msg.SelfID == msg.UserID {
					continue
				}
				if strings.Contains(ignore, fmt.Sprint(msg.UserID)) {
					continue
				}
				if msg.GroupID != 0 {
					if onGroups := qq.Get("offGroups", "923993867"); onGroups != "" && strings.Contains(onGroups, fmt.Sprint(msg.GroupID)) {
						continue
					}
					if onGroups := qq.Get("onGroups"); onGroups != "" && !strings.Contains(onGroups, fmt.Sprint(msg.GroupID)) {
						continue
					}
				}
				// if msg.PostType == "message" {
				msg.RawMessage = strings.ReplaceAll(msg.RawMessage, "\\r", "\n")
				msg.RawMessage = regexp.MustCompile(`[\n\r]+`).ReplaceAllString(msg.RawMessage, "\n")
				core.Senders <- &Sender{
					Conn:    ws,
					Message: msg,
				}
				// }
			}
		}()

	})
}

type Sender struct {
	botID    string
	Conn     *websocket.Conn
	Message  *Message
	matches  [][]string
	Duration *time.Duration
	deleted  bool
	core.BaseSender
}

func (sender *Sender) GetContent() string {
	if sender.Content != "" {
		return sender.Content
	}
	text := sender.Message.RawMessage
	text = strings.Replace(text, "amp;", "", -1)
	text = strings.Replace(text, "&#91;", "[", -1)
	text = strings.Replace(text, "&#93;", "]", -1)

	return strings.Trim(text, " ")
}

func (sender *Sender) GetUserID() string {
	return fmt.Sprint(sender.Message.UserID)
}

func (sender *Sender) GetChatID() int {
	return sender.Message.GroupID
}

func (sender *Sender) GetImType() string {
	return "qq"
}

func (sender *Sender) GetMessageID() int {
	return sender.Message.MessageID
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
	if sender.Message.UserID == sender.Message.SelfID {
		return true
	}
	uid := fmt.Sprint(sender.Message.UserID)
	for _, v := range regexp.MustCompile(`\d+`).FindAllString(qq.Get("masters"), -1) {
		if uid == v {
			return true
		}
	}
	return false
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) GroupKick(uid string, reject_add_request bool) {
	sender.Conn.WriteJSON(CallApi{
		Action: "set_group_kick",
		Params: map[string]interface{}{
			"group_id":           sender.Message.GroupID,
			"user_id":            core.Int(uid),
			"reject_add_request": reject_add_request,
		},
	})
}

func (sender *Sender) GroupBan(uid string, duration int) {
	sender.Conn.WriteJSON(CallApi{
		Action: "set_group_ban",
		Params: map[string]interface{}{
			"group_id": sender.Message.GroupID,
			"user_id":  core.Int(uid),
			"duration": duration,
		},
	})
}

var dd sync.Map

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	chatId := sender.GetChatID()
	if chatId != 0 {
		if onGroups := qq.Get("spy_on", "9251251"); onGroups != "" && strings.Contains(onGroups, fmt.Sprint(chatId)) {
			return 0, nil
		}
	}
	msg := msgs[0]
	rt := ""
	for _, item := range msgs {
		switch item.(type) {
		case time.Duration:
			du := item.(time.Duration)
			sender.Duration = &du
		case string:
			rt = msg.(string)
		case []byte:
			rt = string(msg.([]byte))
		case core.ImageUrl:
			rt = `[CQ:image,file=` + string(msg.(core.ImageUrl)) + `]`
		case core.VideoUrl:
			rt = `[CQ:video,file=` + string(msg.(core.VideoUrl)) + `]`

		}
	}
	if rt == "" {
		return 0, nil
	}
	if sender.Message.MessageType == "private" {
		sender.Conn.WriteJSON(CallApi{
			Action: "send_private_msg",
			Params: map[string]interface{}{
				"user_id": sender.Message.UserID,
				"message": rt,
			},
		})
	} else {
		sender.Conn.WriteJSON(CallApi{
			Action: "send_group_msg",
			Params: map[string]interface{}{
				"group_id": sender.Message.GroupID,
				"user_id":  sender.Message.UserID,
				"message":  rt,
			},
		})
	}
	return 0, nil
}

func (sender *Sender) Delete() error {
	return sender.Conn.WriteJSON(CallApi{
		Action: "delete_msg",
		Params: map[string]interface{}{
			"message_id": sender.Message.MessageID,
		},
	})
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {

}

func (sender *Sender) Copy() core.Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(Sender)
	return &new
}

func (sender *Sender) GetUsername() string {
	return sender.Message.Sender.Nickname
}

func (sender *Sender) GetChatname() string {
	return ""
}
