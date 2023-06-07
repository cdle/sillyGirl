package qq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var qq = core.MakeBucket("qq")

type Result struct {
	Data struct {
		MessageID interface{} `json:"message_id"`
	} `json:"data"`
	Echo string `json:"echo"`
}

type GroupList struct {
	Retcode int    `json:"retcode"`
	Status  string `json:"status"`
	Data    []struct {
		GroupID   int    `json:"group_id"`
		GroupName string `json:"group_name"`
		// MemberCount     int         `json:"member_count"`
		// MaxMemberCount  int         `json:"max_member_count"`
		// OwnerID         int         `json:"owner_id"`
		// LastJoinTime    int         `json:"last_join_time"`
		// ShutupTimeWhole int         `json:"shutup_time_whole"`
		// ShutupTimeMe    int         `json:"shutup_time_me"`
		// AdminFlag       bool        `json:"admin_flag"`
		// UpdateTime      int         `json:"update_time"`
	} `json:"data"`
	Error interface{} `json:"error"`
	Echo  string      `json:"echo"`
}

type CallApi struct {
	Action string                 `json:"action"`
	Echo   string                 `json:"echo"`
	Params map[string]interface{} `json:"params"`
}

type sender struct {
	Nickname string `json:"nickname"`
}

type Message struct {
	GroupID     interface{} `json:"group_id"`
	Message     string      `json:"message"`
	MessageID   interface{} `json:"message_id"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	RawMessage  string      `json:"raw_message"`
	SelfID      interface{} `json:"self_id"`
	Sender      sender      `json:"sender"`
	// SubType     string      `json:"sub_type"`
	Time   int         `json:"time"`
	UserID interface{} `json:"user_id"`
}

type GroupInfo struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type QQ struct {
	conn *websocket.Conn
	sync.RWMutex
	id    int
	chans map[string]chan string
}

func (qq *QQ) WriteJSON(ca CallApi) (string, error) {
	var err error
	cy := make(chan string, 1)
	defer close(cy)
	func() {
		qq.Lock()
		defer qq.Unlock()
		qq.id++
		ca.Echo = fmt.Sprint(qq.id)
		qq.chans[ca.Echo] = cy
		err = qq.conn.WriteJSON(ca)
	}()
	if err != nil {
		return "", err
	}
	select {
	case v := <-cy:
		return v, nil
	case <-time.After(time.Second * 60):
	}
	return "", nil
}

var debug = qq.GetBool("debug", false)

func init() {
	storage.Watch(qq, "debug", func(old, new, key string) *storage.Final {
		now := ""
		if new == "true" {
			now = "true"
			debug = true
		} else {
			now = "false"
			debug = false
		}
		return &storage.Final{
			Now: now,
		}
	})
	go func() {
		core.GinApi(core.GET, "/qq/receive", func(c *gin.Context) {
			auth := c.GetHeader("Authorization")
			token := qq.GetString("access_token")
			if token == "" {
				token = qq.GetString("token")
			}
			if token != "" && !strings.Contains(auth, token) {
				core.Logs.Warn("Onebot机器人access_token不正确，小心有人攻击你的傻妞！！！%s ? %s", token, auth)
			}
			if token == "" {
				core.Logs.Warn(`你需要在Onebot机器人配置access_token以及在傻妞配置对应的参数(set qq access_token ?)才能保证连接安全，如果不设置将会造成信息泄露和资产损失！！！`)
			}
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
			qqcon := &QQ{
				conn:  ws,
				chans: make(map[string]chan string),
			}
			adapter := &core.Factory{}
			adapter.Init("qq", botID)
			defer adapter.Destroy()
			adapter.SetGroupKick(func(uid string, gid string, reject_add_request bool) bool {
				qqcon.WriteJSON(CallApi{
					Action: "set_group_kick",
					Params: map[string]interface{}{
						"group_id":           utils.Int(gid),
						"user_id":            utils.Int(uid),
						"reject_add_request": reject_add_request,
					},
				})
				return true
			})
			adapter.SetGroupBan(func(uid string, gid string, duration int) bool {
				qqcon.WriteJSON(CallApi{
					Action: "set_group_ban",
					Params: map[string]interface{}{
						"group_id": utils.Int(gid),
						"user_id":  utils.Int(uid),
						"duration": duration,
					},
				})
				return true
			})
			adapter.SetGroupUnban(func(uid string, gid string) bool {
				qqcon.WriteJSON(CallApi{
					Action: "set_group_ban",
					Params: map[string]interface{}{
						"group_id": utils.Int(gid),
						"user_id":  utils.Int(uid),
						"duration": 1,
					},
				})
				return true
			})
			adapter.SetReplyHandler(func(msg map[string]string) string {
				if utils.IsZeroOrEmpty(msg[core.CHAT_ID]) {
					id, err := qqcon.WriteJSON(CallApi{
						Action: "send_private_msg",
						Params: map[string]interface{}{
							"user_id": msg[core.USER_ID],
							"message": msg[core.CONETNT],
						},
					})
					if err != nil {
						core.Logs.Warn("QQ发送私聊消息错误：", err)
					}
					return id
				} else {
					id, err := qqcon.WriteJSON(CallApi{
						Action: "send_group_msg",
						Params: map[string]interface{}{
							"group_id": msg[core.CHAT_ID],
							"user_id":  msg[core.USER_ID],
							"message":  msg[core.CONETNT],
						},
					})
					if err != nil {
						core.Logs.Warn("QQ发送群组消息错误：", err)
					}
					return id
				}
			})

			// qqcon.WriteJSON(CallApi{
			// 	Action: "get_group_list",
			// 	Params: map[string]interface{}{},
			// })
			go func() {
				time.Sleep(time.Second * 3)
				qqcon.WriteJSON(CallApi{
					Action: "get_group_list",
					Params: map[string]interface{}{},
				})
			}()

			for {
				_, data, err := ws.ReadMessage()
				if err != nil {
					ws.Close()
					break
				}
				if debug {
					logs.Debug("QQ接收消息：", string(data))
				}
				{
					res := &GroupList{}
					json.Unmarshal(data, res)
					for _, group := range res.Data {
						core.CreateNickName(&core.Nickname{
							Group:    true,
							Value:    group.GroupName,
							ID:       strconv.Itoa(group.GroupID),
							Platform: "qq",
							BotsID:   []string{botID},
						})
					}
				}
				{
					res := &Result{}
					json.Unmarshal(data, res)
					if res.Echo != "" {
						func() {
							qqcon.RLock()
							defer qqcon.RUnlock()
							if c, ok := qqcon.chans[res.Echo]; ok {
								c <- fmt.Sprint(res.Data.MessageID)
							}
						}()
						continue
					}
				}

				msg := &Message{}
				json.Unmarshal(data, msg)
				// if msg.MessageType != "private"} //&& adapter.IsAdapter(botID) {
				// 	continue
				// }
				if msg.SelfID == msg.UserID {
					continue
				}
				msg.RawMessage = strings.ReplaceAll(msg.RawMessage, "\\r", "\n")
				msg.RawMessage = regexp.MustCompile(`[\n\r]+`).ReplaceAllString(msg.RawMessage, "\n")
				content := msg.RawMessage
				content = strings.Replace(content, "amp;", "", -1)
				content = strings.Replace(content, "&#91;", "[", -1)
				content = strings.Replace(content, "&#93;", "]", -1)
				content = strings.Trim(content, " ")
				_msg := map[string]interface{}{
					"user_id":    fmt.Sprint(msg.UserID),
					"chat_id":    core.ChatID(msg.GroupID),
					"user_name":  msg.Sender.Nickname,
					"chat_name":  "",
					"message_id": fmt.Sprint(msg.MessageID),
					"content":    content,
				}
				if debug {
					logs.Debug("QQ处理消息：", string(utils.JsonMarshal(_msg)))
				}
				adapter.Receive(msg)
			}
		})
	}()
}
