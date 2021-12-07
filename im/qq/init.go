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

type Params struct {
	UserID  interface{} `json:"user_id"`
	Message string      `json:"message"`
	GroupID int         `json:"group_id"`
}

type CallApi struct {
	Action string `json:"action"`
	Echo   string `json:"echo"`
	Params Params `json:"params"`
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

var conn *websocket.Conn

func init() {
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
		conn = ws
		core.Pushs["qq"] = func(i interface{}, s string, _ interface{}) {
			if conn == nil {
				return
			}
			conn.WriteJSON(CallApi{
				Action: "send_private_msg",
				Params: Params{
					UserID:  core.Int64(i),
					Message: s,
				},
			})
		}
		core.GroupPushs["qq"] = func(i, j interface{}, s string) {
			if conn == nil {
				return
			}
			conn.WriteJSON(CallApi{
				Action: "send_group_msg",
				Params: Params{
					GroupID: core.Int(i),
					UserID:  core.Int64(j),
					Message: s,
				},
			})
		}
		go func() {
			for {
				_, data, err := ws.ReadMessage()
				if err != nil {
					ws.Close()
					conn = nil
					break
				}
				// fmt.Println(string(data))
				msg := &Message{}
				json.Unmarshal(data, msg)
				// fmt.Println(msg)
				if msg.PostType == "message" {
					msg.RawMessage = strings.ReplaceAll(msg.RawMessage, "\\r", "\n")
					msg.RawMessage = strings.ReplaceAll(msg.RawMessage, "\r", "\n")
					core.Senders <- &Sender{
						Conn:    ws,
						Message: msg,
					}
				}
			}
		}()

	})
}

type Sender struct {
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

var dd sync.Map

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	msg := msgs[0]
	for _, item := range msgs {
		switch item.(type) {
		case time.Duration:
			du := item.(time.Duration)
			sender.Duration = &du
		}
	}
	fmt.Println(msg)
	if sender.Message.MessageType == "private" {
		sender.Conn.WriteJSON(CallApi{
			Action: "send_private_msg",
			Params: Params{
				UserID:  sender.Message.UserID,
				Message: fmt.Sprint(msg),
			},
		})
	} else {
		sender.Conn.WriteJSON(CallApi{
			Action: "send_group_msg",
			Params: Params{
				GroupID: sender.Message.GroupID,
				UserID:  sender.Message.UserID,
				Message: fmt.Sprint(msg),
			},
		})
	}
	return 0, nil
}

func (sender *Sender) Delete() error {

	return nil
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
