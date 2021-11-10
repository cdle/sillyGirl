package wx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var wx = core.NewBucket("wx")
var api_url = func() string {
	return wx.Get("api_url")
}
var robot_wxid = wx.Get("robot_wxid")

func enableBGM() {
	core.Server.POST("/wx/receive", func(c *gin.Context) {
		data, _ := c.GetRawData()
		jms := JsonMsg{}
		json.Unmarshal(data, &jms)
		c.JSON(200, map[string]string{"code": "-1"})
		fmt.Println(jms.Type, jms.Msg)
		if jms.Event != "EventFriendMsg" && jms.Event != "EventGroupMsg" {
			return
		}

		if jms.Type == 0 { //|| jms.Type == 49
			// if jms.Type != 1 && jms.Type != 3 && jms.Type != 5 {
			return
		}
		if strings.Contains(fmt.Sprint(jms.Msg), `<type>57</type>`) {
			return
		}
		if jms.FinalFromWxid == jms.RobotWxid {
			return
		}
		listen := wx.Get("onGroups")
		if jms.Event == "EventGroupMsg" && listen != "" && !strings.Contains(listen, strings.Replace(fmt.Sprint(jms.FromWxid), "@chatroom", "", -1)) {
			return
		}
		if robot_wxid != jms.RobotWxid {
			robot_wxid = jms.RobotWxid
			wx.Set("robot_wxid", robot_wxid)
		}
		if wx.GetBool("keaimao_dynamic_ip", false) {
			ip, _ := c.RemoteIP()
			wx.Set("api_url", fmt.Sprintf("http://%s:%s", ip.String(), wx.Get("keaimao_port", "8080"))) //
		}
		wm := wxmsg{}
		switch jms.Msg.(type) {
		case int, int64, int32:
			wm.content = fmt.Sprintf("%d", jms.Msg)
		case float64:
			wm.content = fmt.Sprintf("%d", int(jms.Msg.(float64)))
		default:
			wm.content = fmt.Sprint(jms.Msg)
		}
		wm.user_id = jms.FinalFromWxid
		wm.user_name = jms.FinalFromName
		if strings.Contains(jms.FromWxid, "@chatroom") {
			wm.chat_id = core.Int(strings.Replace(jms.FromWxid, "@chatroom", "", -1))
		}
		core.Senders <- &Sender{
			value: wm,
		}
	})
}

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 && c != 10 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}
	return string(dstRunes)
}

type TextMsg struct {
	Event      string `json:"event"`
	ToWxid     string `json:"to_wxid"`
	Msg        string `json:"msg"`
	RobotWxid  string `json:"robot_wxid"`
	GroupWxid  string `json:"group_wxid"`
	MemberWxid string `json:"member_wxid"`
}

type OtherMsg struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Event      string `json:"event"`
	RobotWxid  string `json:"robot_wxid"`
	ToWxid     string `json:"to_wxid"`
	MemberWxid string `json:"member_wxid"`
	MemberName string `json:"member_name"`
	GroupWxid  string `json:"group_wxid"`
	Msg        Msg    `json:"msg"`
}

type Msg struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type wxmsg struct {
	content   string
	user_id   string
	chat_id   int
	user_name string
}

type Sender struct {
	leixing int
	mtype   int
	deleted bool
	value   wxmsg
	core.BaseSender
}

type JsonMsg struct {
	Event         string      `json:"event"`
	RobotWxid     string      `json:"robot_wxid"`
	RobotName     string      `json:"robot_name"`
	Type          int         `json:"type"`
	FromWxid      string      `json:"from_wxid"`
	FromName      string      `json:"from_name"`
	FinalFromWxid string      `json:"final_from_wxid"`
	FinalFromName string      `json:"final_from_name"`
	ToWxid        string      `json:"to_wxid"`
	Msg           interface{} `json:"msg"`
}
