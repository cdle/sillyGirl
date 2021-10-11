package wx

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var wx = core.NewBucket("wx")
var api_url = wx.Get("api_url")
var robot_wxid = wx.Get("robot_wxid")

func sendTextMsg(pmsg *TextMsg) {
	pmsg.Msg = TrimHiddenCharacter(pmsg.Msg)
	if pmsg.Msg == "" {
		return
	}
	pmsg.Event = "SendTextMsg"
	pmsg.RobotWxid = robot_wxid
	req := httplib.Post(api_url)
	req.Header("Content-Type", "application/json")
	data, _ := json.Marshal(pmsg)
	enc := mahonia.NewEncoder("gbk")
	d := enc.ConvertString(string(data))
	d = regexp.MustCompile(`[\n\s]*\n[\s\n]*`).ReplaceAllString(d, "\n")
	req.Body(d)
	req.Response()
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

func sendOtherMsg(pmsg *OtherMsg) {
	if pmsg.Event == "" {
		pmsg.Event = "SendImageMsg"
	}
	pmsg.RobotWxid = robot_wxid
	req := httplib.Post(api_url)
	req.Header("Content-Type", "application/json")
	data, _ := json.Marshal(pmsg)
	req.Body(data)
	req.Response()
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

func init() {
	core.Pushs["wx"] = func(i interface{}, s string) {
		if robot_wxid != "" {
			pmsg := TextMsg{
				Msg:    s,
				ToWxid: fmt.Sprint(i),
			}
			sendTextMsg(&pmsg)
		}
	}
	core.GroupPushs["wx"] = func(i, j interface{}, s string) {
		to := fmt.Sprint(i) + "@chatroom"
		pmsg := TextMsg{
			ToWxid: to,
		}
		if j != nil && fmt.Sprint(j) != "" {
			pmsg.MemberWxid = fmt.Sprint(j)
		}
		for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(s, -1) {
			s = strings.Replace(s, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
			data, err := os.ReadFile(core.ExecPath + "/data/images/" + v[1])
			if err == nil {
				add := regexp.MustCompile("(https.*)").FindString(string(data))
				if add != "" {
					pmsg := OtherMsg{
						ToWxid: to,
						Msg: Msg{
							URL:  relay(add),
							Name: name(add),
						},
					}
					defer sendOtherMsg(&pmsg)
				}
			}
		}
		s = regexp.MustCompile(`\[CQ:([^\[\]]+)\]`).ReplaceAllString(s, "")
		pmsg.Msg = s
		sendTextMsg(&pmsg)
	}
	core.Server.POST("/wx/receive", func(c *gin.Context) {
		data, _ := c.GetRawData()
		jms := JsonMsg{}
		json.Unmarshal(data, &jms)
		c.JSON(200, map[string]string{"code": "-1"})
		if jms.Event != "EventFriendMsg" && jms.Event != "EventGroupMsg" {
			return
		}
		if jms.Type != 1 && jms.Type != 3 {
			return
		}
		if jms.FinalFromWxid == jms.RobotWxid {
			return
		}
		if robot_wxid != jms.RobotWxid {
			robot_wxid = jms.RobotWxid
			wx.Set("robot_wxid", robot_wxid)
		}
		core.Senders <- &Sender{
			value: jms,
		}
	})
	core.Server.GET("/relay", func(c *gin.Context) {
		url := c.Query("url")
		rsp, err := httplib.Get(url).Response()
		if err == nil {
			io.Copy(c.Writer, rsp.Body)
		}
	})
}

var myip = ""
var relaier = wx.Get("relaier")

func relay(url string) string {
	if wx.GetBool("relay_mode", false) == false {
		return url
	}
	if relaier != "" {
		return fmt.Sprintf(relaier, url)
	} else {
		if myip == "" || wx.GetBool("dynamic_ip", false) == true {
			ip, _ := httplib.Get("https://imdraw.com/ip").String()
			if ip != "" {
				myip = ip
			}
		}
		return fmt.Sprintf("http://%s:%s/relay?url=%s", myip, wx.Get("relay_port", core.Bucket("sillyGirl").Get("port")), url) //"8002"
	}
}

type Sender struct {
	leixing int
	mtype   int
	deleted bool
	value   JsonMsg
	core.BaseSender
}

type JsonMsg struct {
	Event         string `json:"event"`
	RobotWxid     string `json:"robot_wxid"`
	RobotName     string `json:"robot_name"`
	Type          int    `json:"type"`
	FromWxid      string `json:"from_wxid"`
	FromName      string `json:"from_name"`
	FinalFromWxid string `json:"final_from_wxid"`
	FinalFromName string `json:"final_from_name"`
	ToWxid        string `json:"to_wxid"`
	Msg           string `json:"msg"`
}

func (sender *Sender) GetContent() string {
	return sender.value.Msg
}
func (sender *Sender) GetUserID() interface{} {
	return sender.value.FinalFromWxid
}
func (sender *Sender) GetChatID() interface{} {
	return strings.Replace(sender.value.FromWxid, "@chatroom", "", -1)
}
func (sender *Sender) GetImType() string {
	return "wx"
}
func (sender *Sender) GetUsername() string {
	return sender.value.FinalFromName
}
func (sender *Sender) GetReplySenderUserID() int {
	if !sender.IsReply() {
		return 0
	}
	return 0
}
func (sender *Sender) IsAdmin() bool {
	return strings.Contains(wx.Get("masters"), fmt.Sprint(sender.GetUserID()))
}
func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	to := sender.value.FromWxid
	at := ""
	if to == "" {
		to = sender.value.FinalFromWxid
	} else {
		at = sender.value.FinalFromWxid
	}
	pmsg := TextMsg{
		ToWxid:     to,
		MemberWxid: at,
	}
	for _, item := range msgs {
		switch item.(type) {
		case string:
			pmsg.Msg = item.(string)
		case []byte:
			pmsg.Msg = string(item.([]byte))
		case core.ImageUrl:
			url := string(item.(core.ImageUrl))
			pmsg := OtherMsg{
				ToWxid:     to,
				MemberWxid: at,
				Msg: Msg{
					URL:  relay(url),
					Name: name(url),
				},
			}
			sendOtherMsg(&pmsg)
		}
	}
	if pmsg.Msg != "" {
		sendTextMsg(&pmsg)
	}
	return 0, nil
}

func name(str string) string {
	pr := "jpg"
	ss := regexp.MustCompile(`\.([A-Za-z0-9]+)$`).FindStringSubmatch(str)
	if len(ss) != 0 {
		pr = ss[1]
	}
	md5 := md5V(str)
	return md5 + "." + pr
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
