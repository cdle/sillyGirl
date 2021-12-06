package wx

import (
	"github.com/cdle/sillyGirl/core"
)

var wx = core.NewBucket("wx")
var api_url = func() string {
	return wx.Get("api_url")
}
var robot_wxid = wx.Get("robot_wxid")

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
	chat_name string
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
