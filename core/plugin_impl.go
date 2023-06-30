package core

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	// "golang.org/x/image/webp"
)

var sleep = func(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func GetScriptNameByUUID(uuid string) string {
	for _, f := range Functions {
		if f.UUID == uuid {
			return fmt.Sprintf("%s.js", f.Title)
		}
	}
	return "未知脚本.js"
}

type SenderJsIplm struct {
	Message    common.Sender
	Private    string
	Group      string
	Routine    string
	Persistent string
	Vm         *goja.Runtime
	UUID       string
	NewContent bool
}

type Console struct {
	UUID string
}

var console = &Console{}
var Logs = &Console{}

func Broadcast2WebUser(content, class string) {
	if (RegistFuncs["Broadcast2WebUser"]) == nil {
		return
	}
	RegistFuncs["Broadcast2WebUser"].(func(string, string))(content, class)
}

func (c *Console) Info(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	log := utils.FormatLog(v[0], v[1:]...)
	logs.Info(log)
	Broadcast2WebUser(log, "info")
}

func (c *Console) Debug(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	log := utils.FormatLog(v[0], v[1:]...)
	logs.Debug(log)
	Broadcast2WebUser(log, "debug")
}

func (c *Console) Warn(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	log := utils.FormatLog(v[0], v[1:]...)
	logs.Warn(log)
	WritePluginMessage(c.UUID, "warn", log)
	Broadcast2WebUser(log, "warn")
}

func (c *Console) Error(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	log := utils.FormatLog(v[0], v[1:]...)
	logs.Error(log)
	WritePluginMessage(c.UUID, "error", log)
	Broadcast2WebUser(log, "error")
}

func (c *Console) Log(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	log := utils.FormatLog(v[0], v[1:]...)
	logs.Info(log)
	Broadcast2WebUser(log, "log")
}

func (sender *SenderJsIplm) Continue() {
	if sender.NewContent {
		go func() {
			Messages <- sender.Message
		}()
		return
	}
	sender.Message.Continue()
}

func (sender *SenderJsIplm) GetUserID() string {

	return sender.Message.GetUserID()
}
func (sender *SenderJsIplm) GetBotID() string {

	return sender.Message.GetBotID()
}
func (sender *SenderJsIplm) GetBotId() string {

	return sender.Message.GetBotID()
}

func (sender *SenderJsIplm) GetUserId() string {
	return sender.Message.GetUserID()
}

func (sender *SenderJsIplm) SetContent(s string) {
	if s != sender.GetContent() {
		sender.NewContent = true
	}
	sender.Message.SetContent(s)
}

func (sender *SenderJsIplm) GetContent() string {
	return sender.Message.GetContent()
}
func (sender *SenderJsIplm) GetImType() string {
	return sender.Message.GetImType()
}
func (sender *SenderJsIplm) GetPlatform() string {
	return sender.Message.GetImType()
}
func (sender *SenderJsIplm) RecallMessage(p ...interface{}) {
	np := []interface{}{}
	var i = 0
	for _, v := range p {
		switch v.(type) {
		case int, int32, int64, uint, int16:
			i = utils.Int(v)
		default:
			np = append(np, v)
		}
	}
	if i != 0 {
		go func() {
			sleep(i)
			sender.Message.RecallMessage(np...)
		}()
	} else {
		go sender.Message.RecallMessage(np...)
	}
}
func (sender *SenderJsIplm) GetUserName() string {
	return sender.Message.GetUserName()
}
func (sender *SenderJsIplm) GetUsername() string {
	return sender.Message.GetUserName()
}

func (sender *SenderJsIplm) GetReplyUserID() int {
	return sender.Message.GetReplyUserID()
}

func (sender *SenderJsIplm) GetReplyUserId() int {
	return sender.Message.GetReplyUserID()
}

func (sender *SenderJsIplm) GetLevel() int {
	return sender.Message.GetLevel()
}

func (sender *SenderJsIplm) SetLevel(l int) {
	sender.Message.SetLevel(l)
}

func (sender *SenderJsIplm) GetChatName() string {
	return sender.Message.GetChatName()
}
func (sender *SenderJsIplm) GetMessageID() *goja.Promise {
	promise, resolve, _ := sender.Vm.NewPromise()
	go func() {
		resolve(sender.Message.GetMessageID())
	}()
	return promise
}

func (sender *SenderJsIplm) GetMessageId() string {
	return sender.Message.GetMessageID()
}

func (sender *SenderJsIplm) GetGroupCode() string {
	return sender.Message.GetChatID()
}

func (sender *SenderJsIplm) GetChatID() string {
	return sender.Message.GetChatID()
}

func (sender *SenderJsIplm) GetChatId() string {
	return sender.Message.GetChatID()
}

func (sender *SenderJsIplm) Kick(uid string) {
	sender.Message.GroupKick(uid, false)
}

func (sender *SenderJsIplm) Unkick(uid string) {
	sender.Message.GroupUnkick(uid)
}

func (sender *SenderJsIplm) Ban(uid string, duration int) {
	sender.Message.GroupBan(uid, duration)
}

func (sender *SenderJsIplm) Param(i interface{}) string {
	switch i := i.(type) {
	case int:
		return sender.Message.Get(i - 1)
	case int64:
		return sender.Message.Get(int(i - 1))
	case string:
		return sender.Message.Get(i)
	}
	return ""
}

func (sender *SenderJsIplm) GetAllMatch() [][]string {
	return sender.Message.GetAllMatch()
}

func (sender *SenderJsIplm) Params(i int) []string {
	if i == 0 {
		i = 1
	}
	ss := []string{}
	for _, v := range sender.Message.GetAllMatch() {
		ss = append(ss, v[i-1])
	}
	return ss
}

// GetAllMatch

func (Sender *SenderJsIplm) Get(i interface{}) string {
	return Sender.Param(i)
}

func (Sender *SenderJsIplm) GetVar(key string) interface{} {
	return Sender.Message.GetVar(key)
}

func (Sender *SenderJsIplm) SetVar(key string, value interface{}) {
	Sender.Message.SetVar(key, value)
}

func (Sender *SenderJsIplm) SetVars(kvs map[string]interface{}) {
	for k, v := range kvs {
		Sender.SetVar(k, v)
	}
}

func (Sender *SenderJsIplm) GetVars() map[string]interface{} {
	return Sender.Message.GetExpandMessageInfo()
}

func (sender *SenderJsIplm) IsAdmin() bool {
	return sender.Message.IsAdmin()
}

func (sender *SenderJsIplm) Reply(texts ...interface{}) interface{} {
	i, err := sender.Message.Reply(texts...)
	if err != nil {
		panic(Error(sender.Vm, err))
	}
	return i
}

func (sender *SenderJsIplm) HoldOn(str ...string) string {
	if len(str) != 0 {
		return "go_again_" + str[0]
	}
	return "go_again_"
}

func arrayss(i interface{}) []string {
	var ss = []string{}
	for _, v := range i.([]interface{}) {
		ss = append(ss, v.(string))
	}
	return ss
}

func (sender *SenderJsIplm) Listen(ps ...interface{}) interface{} {
	// promise, resolve, reject := vm.NewPromise()
	options := []interface{}{}
	var handle func(goja.FunctionCall) goja.Value
	var persistent = sender.Message.GetImType() == "*"
	var routine = persistent
	var carry = &Carry{}
	// var sustainable = false
	// rs := []string{}
	for _, p := range ps {
		switch p := p.(type) {
		case map[string]interface{}:
			props := p
			for i, p := range props {
				switch strings.ToLower(i) {
				case "rules":
					vs := p.([]interface{})
					for _, v := range vs {
						rule := v.(string)
						_rs := formatRule(rule)
						if len(_rs) != 0 {
							carry.Function.Rules = append(carry.Function.Rules, _rs...)
						} else {
							carry.Function.Rules = append(carry.Function.Rules, rule)
						}
					}
				case "timeout":
					options = append(options, time.Duration(utils.Int(p))*time.Millisecond)
				case "handle":
					handle = p.(func(goja.FunctionCall) goja.Value)
				case "listen_private", "private":
					carry.ListenPrivate = p.(bool)
				case "listen_group", "group":
					carry.ListenGroup = p.(bool)
				case "require_admin", "admin":
					carry.RequireAdmin = p.(bool)
				case "allow_platforms":
					carry.AllowPlatforms = arrayss(p)
				case "prohibit_platforms":
					carry.ProhibitPlatforms = arrayss(p)
				case "allow_groups":
					carry.AllowGroups = arrayss(p)
				case "allow_users":
					carry.AllowUsers = arrayss(p)
				case "prohibt_groups":
					carry.ProhibitGroups = arrayss(p)
				case "prohibt_users":
					carry.ProhibitUsers = arrayss(p)
				}
			}
		case bool:
		case []interface{}:
			for i := range p {
				rule := p[i].(string)
				_rs := formatRule(rule)
				if len(_rs) != 0 {
					carry.Function.Rules = append(carry.Function.Rules, _rs...)
				} else {
					carry.Function.Rules = append(carry.Function.Rules, rule)
				}
			}
		case string:
			if p == "private" {
				carry.ListenPrivate = true
				continue
			}
			if p == "group" {
				carry.ListenGroup = true
				continue
			}
			if p == "routine" {
				routine = true
				continue
			}
			if p == "persistent" {
				persistent = true
				continue
			}
			// if p == "sustainable" {
			// 	persistent = false
			// 	sustainable = true
			// 	continue
			// }
			carry.Function.Rules = append(carry.Function.Rules, p)
		case int, int64, int32:
			options = append(options, time.Duration(utils.Int(p))*time.Millisecond)
		case func(goja.FunctionCall) goja.Value:
			handle = p
		}
	}
	if persistent {
		options = append(options, "persistent")
	}
	if !persistent {
		carry.AllowPlatforms = []string{sender.GetPlatform()}
	}
	carry.UUID = sender.UUID
	options = append(options, carry)
	var newJsSender *SenderJsIplm
	var await = func() {
		sender.Message.Await(sender.Message, func(newSender common.Sender) interface{} {
			newSender.SetLevel(sender.Message.GetLevel() + 1)
			newJsSender = &SenderJsIplm{
				UUID: sender.UUID,
				Vm:   sender.Vm,
			}
			newJsSender.Message = newSender
			if handle != nil {
				if persistent {
					func() {
						defer func() {
							err := recover()
							if err != nil {
								console.Error("%v at %v", err, GetScriptNameByUUID(sender.UUID))
							}
						}()
						// fmt.Println(newJsSender)
						rt := handle(goja.FunctionCall{
							Arguments: []goja.Value{
								sender.Vm.ToValue(newJsSender),
							},
						})
						reply := ""
						if rt != goja.Undefined() {
							reply = rt.ToString().String()
						} else {
							reply = ""
						}
						newSender.Reply(reply)
					}()
					// return GoAgain("")
				} else {
					var rt = handle(goja.FunctionCall{
						Arguments: []goja.Value{
							sender.Vm.ToValue(newJsSender),
						},
					})
					reply := ""
					if rt != goja.Undefined() {
						reply = rt.ToString().String()
					} else {
						reply = ""
					}
					if strings.HasPrefix(reply, "go_again_") {
						reply = strings.Replace(reply, "go_again_", "", 1)
						return GoAgain(reply)
					} else {
						// if sustainable {
						// 	return GoAgain(reply)
						// }
						if reply == "" {
							return nil
						}
						return reply
					}
				}
			}
			return nil
		}, options...)
	}
	if !routine {
		await()
	} else {
		go await()
	}
	return newJsSender
}

type TimeJsImpl struct {
	Second time.Duration
	Minute time.Duration
	Hour   time.Duration
	Day    time.Duration
	Month  int
}

func (t *TimeJsImpl) Now() time.Time {
	return time.Now()
}

func (t *TimeJsImpl) Date(year int, month int, day int, hour int, min int, sec int, nsec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc)
}

func (t *TimeJsImpl) Sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func (t *TimeJsImpl) Unix(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func (t *TimeJsImpl) UnixMilli(sec int64) time.Time {
	return time.UnixMilli(sec)
}

func (t *TimeJsImpl) Parse(timeStr, layout, locale string) (tt time.Time, err error) {
	if locale == "" {
		locale = "Asia/Shanghai"
	}
	local, err := time.LoadLocation(locale) //设置时区
	if err != nil {
		return tt, err
	}
	tt, err = time.ParseInLocation(layout, timeStr, local)
	if err != nil {
		return tt, err
	}
	return tt, nil
}

type Fmt struct {
}

func (sender *Fmt) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func (sender *Fmt) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (sender *Fmt) Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func (sender *Fmt) Print(a ...interface{}) (int, error) {
	return fmt.Print(a...)
}

func Url2Base64(imageUrl string) map[string]interface{} {
	response, err := http.Get(imageUrl)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	imageBase64 := base64.StdEncoding.EncodeToString(data)
	return map[string]interface{}{"result": imageBase64}
}

func stringToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64ToString(b64 string) string {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return ""
	}
	return string(data)
}

func formatRule(rule string) []string {
	_rs := []string{}
FR:
	ress := regexp.MustCompile(`\[([^\s\[\]]+)\]`).FindAllStringSubmatch(rule, -1)
	if len(ress) != 0 {
		res := ress[len(ress)-1]
		var inner = res[1]
		slice := strings.SplitN(inner, ":", 2)
		name := slice[0]
		ps := ""
		if len(slice) == 2 {
			ps = slice[1]
		}
		if strings.HasSuffix(name, "?") {
			name = strings.TrimRight(name, "?")
			rep := ""
			if ps == "" {
				rep = fmt.Sprintf("[%s]", name)
			} else {
				rep = fmt.Sprintf("[%s:%s]", name, ps)
			}
			for l := range _rs {
				_rs[l] = strings.Replace(_rs[l], res[0], rep, 1)
			}
			rule1 := strings.Replace(rule, res[0], rep, 1)
			if len(_rs) == 0 {
				_rs = append(_rs, rule1)
			}
			rule = strings.Replace(rule, res[0], "", 1)
			rule = regexp.MustCompile("\x20{2,}").ReplaceAllString(rule, " ")
			rule = strings.TrimSpace(rule)
			_rs = append(_rs, rule)
			goto FR
		}
	}
	return _rs
}
