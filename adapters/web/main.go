package web

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

type WebMessage struct {
	UserID  string   `json:"-"`
	Type    string   `json:"t"`
	Content string   `json:"c"`
	Images  []string `json:"m"`
}

var webUsers sync.Map

func Broadcast2WebUser(content, class string) {
	webUsers.Range(func(key, value interface{}) bool {
		wu := value.(*WebUser)
		if wu.ActivedAt.Add(10 * time.Second).After(time.Now()) {

			wu.GetCarry() <- WebMessage{
				UserID:  key.(string),
				Content: content,
				Type:    class,
			}
		} else {
			webUsers.Delete(key)
		}
		return true
	})
}

type WebUser struct {
	Carry chan WebMessage
	sync.RWMutex
	ActivedAt time.Time
}

func (wu *WebUser) GetCarry() chan WebMessage {
	wu.RLock()
	defer wu.RUnlock()
	return wu.Carry
}

func (wu *WebUser) Active() {
	wu.Lock()
	defer wu.Unlock()
	wu.ActivedAt = time.Now()
}

func (wu *WebUser) GetActivedAt() time.Time {
	wu.RLock()
	defer wu.RUnlock()
	return wu.ActivedAt
}

var webAdmins sync.Map

var adapter *core.Factory

var GetUserNumber = func() int {
	i := 0
	webUsers.Range(func(key, value any) bool {
		i++
		return true
	})
	return i
}

func initWebBot() {
	if adapter == nil {
		adapter = &core.Factory{}
		adapter.Init("web", "default")
		adapter.SetIsAdmin(func(s string) bool {
			isAdmin, ok := webAdmins.Load(s)
			if ok {
				return isAdmin.(bool)
			}
			return false
		})
		adapter.SetReplyHandler(func(msg map[string]interface{}) string {
			message := WebMessage{
				UserID:  msg[core.USER_ID].(string),
				Images:  []string{},
				Type:    "chat",
				Content: msg[core.CONETNT].(string),
			}
			sendWebMessage(&message)
			return ""
		})
	}
}

func init() {
	core.RegistFuncs["Broadcast2WebUser"] = Broadcast2WebUser
	go func() {
		time.Sleep(time.Second)
		initWebBot()
	}()
	core.GinApi(core.GET, "/api/web_chat", func(ctx *gin.Context) {
		initWebBot()
		rid := ctx.Query("rid")
		ctt := ctx.Query("ctt")
		token, _ := ctx.Cookie("token")
		_, err := core.CheckAuth(token)
		isAdmin := err == nil
		v, ok := webAdmins.Load(rid)
		if ok {
			if v.(bool) != isAdmin {
				webAdmins.Store(rid, isAdmin)
			}
		} else {
			webAdmins.Store(rid, isAdmin)
		}
		if ctt != "" {
			adapter.Receive(map[string]interface{}{
				core.USER_ID: rid,
				core.CONETNT: ctt,
			})
		}
		var wu *WebUser
		if v, ok := webUsers.Load(rid); ok {
			wu = v.(*WebUser)
		} else {
			wu = &WebUser{
				Carry: make(chan WebMessage, 1000),
			}
			webUsers.Store(rid, wu)
		}
		wu.Active()
		msgs := []WebMessage{}
		for {
			select {
			case msg := <-wu.GetCarry():
				msgs = append(msgs, msg)
			case <-time.After(time.Millisecond * 1):
				goto HELL
			}
		}
	HELL:
		if len(msgs) == 0 && ctt == "" {
			select {
			case msg := <-wu.GetCarry():
				msgs = append(msgs, msg)
			case <-time.After(time.Second * 4):
				break
			}
		}
		ctx.JSON(200, msgs)
	})
}

var sendWebMessage = func(message *WebMessage) {
	message.Content = regexp.MustCompile(`file=[^\[\]]*,url`).ReplaceAllString(message.Content, "file")
	for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(message.Content, -1) {
		message.Images = append(message.Images, v[1])
		message.Content = strings.Replace(message.Content, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
	}
	v, ok := webUsers.Load(message.UserID)
	var wu *WebUser
	if !ok {
		wu = &WebUser{
			Carry: make(chan WebMessage, 1000),
		}
		webUsers.Store(message.UserID, wu)
		wu.Active()
	} else {
		wu = v.(*WebUser)
	}
	go func() {
		wu.GetCarry() <- *message
	}()
}
