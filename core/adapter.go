package core

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
)

type Message map[string]string

type Details struct {
	Content   string
	UserID    string
	ChatID    string
	Username  string
	Chatname  string
	MessageID string
}

type CustomSender struct {
	BaseSender
	details Details
	f       *Factory
}

type MsgChan struct {
	Chan chan string
	Msg  map[string]interface{}
}

type GMsgChan struct {
	Chan       chan map[string]interface{}
	Msgs       []map[string]interface{}
	UpdateTime time.Time
	sync.Mutex
}

type Factory struct {
	botid      string
	botplt     string
	uuid       string
	msgChan    chan MsgChan
	demo       *CustomSender
	reply      func(map[string]interface{}) string
	lm         chan bool
	nm         int64
	isAdmin    func(string) bool
	vm         *goja.Runtime
	ctx        context.Context
	cancel     context.CancelFunc
	destroid   bool ////已关闭
	errorTimes int
	Res        *Response
	umod       bool //类似订阅号一对一被动消息模式
	gmsgChan   sync.Map
	sync.RWMutex
}

type Bot [2]string //botplt botid

var Bots = map[Bot]*Factory{}
var BotsLocker sync.RWMutex

var ErrNotFind = errors.New("adapter not find")

func DestroyAdapterByUUID(uuid string) {
	// BotsLocker.RLock()
	// defer BotsLocker.RUnlock()
	// for i := range Bots {
	// 	if Bots[i].uuid == uuid {
	// 		go Bots[i].Destroy()
	// 	}
	// }
}

func GetMessageByUUID(uuid string) string {
	ss := []string{""}
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	for i, bot := range Bots {
		plt, id := i[0], i[1]
		// fmt.Println("plt", plt, "id", id, botplt, bots_id)
		if bot.uuid == uuid {
			ss[0] = plt
			ss = append(ss, id)
		}
	}
	if len(ss) == 0 {
		return ""
	}
	return strings.Join(ss, " ")
}

func GetAdapter(botplt string, bots_id ...string) (*Factory, error) {
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	var bots = []*Factory{}
	var select_bots = []*Factory{}
	for i := range Bots {
		plt, id := i[0], i[1]
		// fmt.Println("plt", plt, "id", id, botplt, bots_id)
		for j := range bots_id {
			if plt == botplt && bots_id[j] == id {
				select_bots = append(select_bots, Bots[i])
			}
			if plt == botplt {
				bots = append(bots, Bots[i])
			}
		}
	}
	if len(bots) == 0 {
		return nil, ErrNotFind
	}
	if len(select_bots) != 0 {
		i := rand.Intn(len(select_bots))
		return select_bots[i], nil
	}
	i := rand.Intn(len(bots))
	return bots[i], ErrNotFind
}

func GetAdapters(botplt string, bots_id ...string) ([]*Factory, error) {
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	var bots = []*Factory{}
	var select_bots = []*Factory{}
	for i := range Bots {
		plt, id := i[0], i[1]
		// fmt.Println("plt", plt, "id", id, botplt, bots_id)
		for j := range bots_id {
			if plt == botplt && bots_id[j] == id {
				select_bots = append(select_bots, Bots[i])
			}
			if plt == botplt {
				bots = append(bots, Bots[i])
			}
		}
	}
	if len(bots) == 0 {
		return nil, ErrNotFind
	}
	if len(select_bots) != 0 {
		return select_bots, nil
	}
	return bots, ErrNotFind
}

func GetAdapterBotsID(botplt string) []string {
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	var bots_id = []string{}
	for i := range Bots {
		// fmt.Println("==", botplt == i[0], botplt, i[0])
		if botplt == i[0] {
			bots_id = append(bots_id, i[1])
		}
	}
	return bots_id
}
func GetAdapterBotPlts() []string {
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	var bot_plts = []string{}
	for i := range Bots {
		// fmt.Println("==", botplt == i[0], botplt, i[0])
		has := false
		for _, bot_plt := range bot_plts {
			if bot_plt == i[0] {
				has = true
			}
		}
		if !has {
			bot_plts = append(bot_plts, i[0])
		}
	}
	return bot_plts
}

func (f *Factory) Init(botplt, botid string, params map[string]interface{}) {
	// fmt.Println(params)
	if params != nil {
		if _, ok := params["umod"]; ok {
			f.umod = true
		}
	}
	// fmt.Println("umod", f.umod)
	f.ctx, f.cancel = context.WithCancel(context.Background())
	BotsLocker.Lock()
	defer BotsLocker.Unlock()
	f.botplt = botplt
	f.botid = botid
	f.msgChan = make(chan MsgChan, 100000)
	f.demo = &CustomSender{
		f: f,
	}
	if v, ok := Bots[[2]string{botplt, botid}]; ok {
		if v.uuid != f.uuid {
			go v.Destroy()
		}
		//
		console.Warn("%s机器人%s因冲突销毁！", botplt, botid)
	}
	Bots[[2]string{botplt, botid}] = f
	f.lm = make(chan bool, 10)
	f.nm = 0
	if botid != "" {
		botid = fmt.Sprintf("(%s)", botid)
	}
	console.Log("%s机器人%s已初始化", botplt, botid)
	go func() {
		if f.uuid != "" {
			su := &ScriptUtils{
				script: plugins.GetString(f.uuid),
			}
			str := su.GetValue("message")
			ss := regexp.MustCompile(`\S+`).FindAllString(str, -1)
			if len(ss) == 0 {
				ss = []string{f.botplt}
			}
			if ss[0] != f.botplt {
				ss = []string{f.botplt}
			}
			nss := utils.Unique(ss, f.botid)
			nstr := strings.Join(nss, " ")
			if str != nstr {
				su.SetValue("message", nstr)
				plugins.Set(f.uuid, su.script)
			}
		}
	}()
}

func (f *Factory) Fail() int {
	f.errorTimes++
	if f.errorTimes > 5 {
		go f.Destroy()
	}
	return f.errorTimes
}

func (f *Factory) Success() {
	f.errorTimes = 0
}

func (f *Factory) IsAdapter(botid string) bool {
	BotsLocker.RLock()
	defer BotsLocker.RUnlock()
	for i := range Bots {
		id := i[0]
		if botid == id {
			return true
		}
	}
	return false
}

func (f *Factory) Masters() []string {
	return strings.Split(strings.Trim(MakeBucket(f.botplt).GetString("masters"), "&"), "&")
}

func (f *Factory) Destroy() {
	BotsLocker.Lock()
	defer BotsLocker.Unlock()
	f.Lock()
	defer f.Unlock()
	if f.destroid {
		return
	}
	f.destroid = true
	f.cancel()
	close(f.msgChan)
	delete(Bots, [2]string{f.botplt, f.botid})
	botid := ""
	if f.botid != "" {
		botid = fmt.Sprintf("(%s)", f.botid)
	} else {
		botid = f.botid
	}
	console.Log("%s机器人%s已销毁", strings.ToUpper(f.botplt), botid)
	go func() {
		if f.uuid != "" {
			su := &ScriptUtils{
				script: plugins.GetString(f.uuid),
			}
			str := su.GetValue("message")
			if str == "" {
				return
			}
			ss := regexp.MustCompile(`\S+`).FindAllString(str, -1)
			if len(ss) == 0 {
				return
			}
			if ss[0] != f.botplt {
				ss = []string{f.botplt}
			}
			nss := utils.Unique(utils.Remove(ss, f.botid))
			if len(nss) == 1 {
				su.DeleteValue("message")
				plugins.Set(f.uuid, su.script)
			} else {
				nstr := strings.Join(nss, " ")
				if str != nstr {
					su.SetValue("message", nstr)
					plugins.Set(f.uuid, su.script)
				}
			}
		}
	}()
}

func (f *Factory) Push(msg map[string]string) (string, error) {
	var demo = *f.demo
	var sender = &demo
	fsps := &common.FakerSenderParams{
		UserID: msg[USER_ID],
		ChatID: msg[CHAT_ID],
	}
	sender.SetFsps(fsps)
	return sender.Reply(msg[CONETNT], PUSH(""))
}

func (f *Factory) SetReplyHandler(function func(map[string]interface{}) string) {
	f.reply = func(m map[string]interface{}) string {
		if f.uuid != "" {
			mutex := GetMutex(f.uuid)
			mutex.Lock()
			defer mutex.Unlock()
		}
		defer func() {
			err := recover()
			if err != nil {
				pluginConsole(f.uuid).Error("Sender(\""+f.botplt+"\").SetReply error:", err)
			}
		}()
		return function(m)
	}
}

// func (f *Factory) GetReplies() {

// }

func GetReplyMessage(vm *goja.Runtime, plt string, bots_id []string) *goja.Promise {
	promise, resolve, reject := vm.NewPromise()
	adapters, err := GetAdapters(plt, bots_id...)
	if err != nil {
		go func() {
			time.Sleep(time.Second)
			reject(Error(vm, err))
		}()
		return promise
	}
	ctx, cancel := context.WithCancel(context.Background())
	for i := range adapters {
		go func(i int) {
			select {
			case <-ctx.Done():
				logs.Debug("消息获取中断")
			case <-adapters[i].ctx.Done():
				cancel()
				logs.Debug("%s adapter %s destroied", adapters[i].botplt, adapters[i].botid)
				reject("adapter destroied")
			case mc := <-adapters[i].msgChan:
				cancel()
				var msg = map[string]interface{}{}
				for k, v := range mc.Msg {
					msg[k] = v
				}
				obj := adapters[i].vm.NewObject()
				for k, v := range mc.Msg {
					obj.Set(k, v)
				}
				obj.Set("bot_id", adapters[i].botid)
				// msg["bot_id"] = adapters[i].botid
				// msg["setMessageId"] = func(id string) {
				// 	select {
				// 	case <-mc.Chan:
				// 	case <-time.After(time.Millisecond):
				// 		mc.Chan <- id
				// 	}
				// }
				// resolve(msg)
				resolve(vm.NewProxy(obj, &goja.ProxyTrapConfig{
					Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
						if property == "message_id" {
							select {
							case <-mc.Chan:
								return false
							case <-time.After(time.Millisecond):
							}
							mc.Chan <- fmt.Sprint(value.Export())
						}
						return true
					},
				}))
			}
		}(i)
	}
	return promise
}

func (f *Factory) GetReplyMessage() interface{} {
	select {
	case <-f.ctx.Done():
		logs.Debug("%s adapter %s destroied", f.botplt, f.botid)
		panic(Error(f.vm, "adapter destroied"))
	case mc := <-f.msgChan:
		obj := f.vm.NewObject()
		for k, v := range mc.Msg {
			obj.Set(k, v)
		}
		return f.vm.NewProxy(obj, &goja.ProxyTrapConfig{
			Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
				if property == "message_id" {
					select {
					case <-mc.Chan:
						return false
					case <-time.After(time.Millisecond):
					}
					mc.Chan <- fmt.Sprint(value.Export())
				}
				return true
			},
		})
	}
}

func (f *Factory) GetUserMessages(user_id string, timeout int) []map[string]interface{} {
	if timeout == 0 {
		timeout = 2000
	}
	msgs := []map[string]interface{}{}
	v, loaded := f.gmsgChan.LoadOrStore(user_id, &GMsgChan{})
	ch := v.(*GMsgChan)
	if !loaded {
		// console.Debug("接收创建：", ch.Chan)
		ch.Chan = make(chan map[string]interface{})
	} else {
		// console.Debug("接收加载：", ch.Chan)
	}
	if len(ch.Msgs) != 0 {
		ch.Lock()
		msgs = append(msgs, ch.Msgs...)
		// console.Debug("数组接收：", msgs)
		ch.Msgs = nil
		ch.Unlock()
		timeout = 1
	}
	for {
		select {
		case msg := <-ch.Chan:
			msgs = append(msgs, msg)
			timeout = 1
			// console.Debug("管道接收：", msgs)
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			// console.Debug("无消息")
			goto HELL
		}
	}
HELL:
	return msgs
}

func (f *Factory) GetMessages(timeout int) []interface{} {
	if timeout == 0 {
		timeout = 2000
	}
	msgs := []interface{}{}
	for {
		select {
		case <-f.ctx.Done():
			logs.Debug("%s adapter %s destroied", f.botplt, f.botid)
			panic(Error(f.vm, "adapter destroied"))
		case mc := <-f.msgChan:
			obj := f.vm.NewObject()
			for k, v := range mc.Msg {
				obj.Set(k, v)
			}
			msgs = append(msgs, f.vm.NewProxy(obj, &goja.ProxyTrapConfig{
				Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
					if property == "message_id" {
						select {
						case <-mc.Chan:
							return false
						case <-time.After(time.Millisecond):
						}
						mc.Chan <- fmt.Sprint(value.Export())
					}
					return true
				},
			}))
			timeout = 1
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			goto HELL
		}
	}
HELL:
	return msgs
}

func (f *Factory) Send(function func(map[string]interface{}) string) {
	f.SetReplyHandler(function)
}

func (f *Factory) SetIsAdmin(function func(string) bool) {
	f.isAdmin = func(uid string) bool {
		defer func() {
			err := recover()
			if err != nil {
				pluginConsole(f.uuid).Error("Sender(\""+f.botplt+"\").SetAdmin error:", err)
			}
		}()
		return function(uid)
	}
}

func (f *Factory) Sender() *CustomSender {
	var demo = *f.demo
	sender := &demo
	return sender
}

func (f *Factory) Receive(wt interface{}) *CustomSender {
	var demo = *f.demo
	sender := &demo
	props := wt.(map[string]interface{})
	emf := map[string]interface{}{}
	for i := range props {
		h := false
		switch strings.ToLower(i) {
		case "content":
			sender.details.Content = strings.TrimSpace(fmt.Sprint(props[i]))
			h = true
		case "message_id", "messageId":
			sender.details.MessageID = utils.Itoa(props[i])
			h = true
		case "user_id", "userId":
			sender.details.UserID = utils.Itoa(props[i])
			h = true
		case "chat_id", "chatId", "group_id", "groupId", "group_code", "groupCode":
			sender.details.ChatID = ChatID(props[i])
			h = true
		case "user_name", "userName":
			sender.details.Username = fmt.Sprint(props[i])
			h = true
		case "chat_name", "chatName", "groupName", "group_name":
			sender.details.Chatname = fmt.Sprint(props[i])
			h = true
		}
		if !h {
			emf[i] = props[i]
		}
	}
	sender.SetExpandMessageInfo(emf)
	// if sender.details.Content != "" {
	Messages <- sender
	// }
	return sender
}

func (sender *CustomSender) GetContent() string {
	if sender.Fsps.Content != "" {
		return sender.Fsps.Content
	}
	return sender.details.Content
}
func (sender *CustomSender) GetUserID() string {
	if sender.Fsps.UserID != "" {
		return sender.Fsps.UserID
	}
	return sender.details.UserID
}
func (sender *CustomSender) GetChatID() string {
	if !utils.IsZeroOrEmpty(sender.Fsps.ChatID) {
		return fmt.Sprint(sender.Fsps.ChatID)
	}
	return sender.details.ChatID
}
func (sender *CustomSender) GetImType() string {
	return sender.f.botplt
}
func (sender *CustomSender) GetUserName() string {
	return sender.details.Username
}
func (sender *CustomSender) GetChatName() string {
	return sender.details.Chatname
}
func (sender *CustomSender) GetMessageID() string {
	return sender.details.MessageID
}
func (sender *CustomSender) GetReplySenderUserID() int {
	if !sender.IsReply() {
		return 0
	}
	return 0
}

func (sender *CustomSender) GetBotID() string {
	return sender.f.botid
}

type PUSH string

func (sender *CustomSender) Action(options map[string]interface{}) (interface{}, string) {
	var platform = sender.f.botplt
	var any *common.Function
	var one *common.Function
	var result interface{}
	var err = ""
	for _, function := range Functions {
		if function.Reply != nil && function.Reply.Platform == platform {
			if len(function.Reply.BotsID) == 0 {
				any = function
			} else if Contains(function.Reply.BotsID, sender.f.botid) {
				one = function
			}
		}
	}
	if one == nil && any != nil {
		one = any
	}
	if one != nil {
		one.Handle(&Faker{
			Type: "action",
		}, func(vm *goja.Runtime) {
			obj := vm.NewObject()
			for k, v := range options {
				obj.Set(k, v)
			}
			proxy := vm.NewProxy(obj, &goja.ProxyTrapConfig{
				Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
					if property == "result" {
						result = value
					}
					if property == "error" {
						err = value.String()
					}
					return true
				},
			})
			vm.Set("action", proxy)
			vm.Set("adapter", sender.f)
		})
	}
	return result, err
}

func (sender *CustomSender) Reply(msgs ...interface{}) (string, error) {
	var push = false
	var platform = sender.f.botplt
	var bot_id = sender.f.botid
	var args = []interface{}{}
	for _, item := range msgs {
		switch item := item.(type) {
		case PUSH:
			push = true
		case string:
			args = append(args, item)
		default:
			if item != nil {
				args = append(args, fmt.Sprint(item))
			}
		}
	}
	if len(args) == 0 {
		return "", nil //errors.New("s.reply has no content")
	}
	content := utils.FormatLog(args[0], args[1:]...)
	if !push {
		if IsNoReplyGroup(sender) {
			return "", errors.New("is no reply group")
		}
	}
	content = strings.ReplaceAll(content, "\n\r", "\n")
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = regexp.MustCompile("[\n]{3,}").ReplaceAllString(content, "\n\n")
	if content != "" {
		user_id := sender.GetUserID()
		msg := map[string]interface{}{
			"message_id": sender.GetMessageID(),
			"content":    content,
			"user_id":    user_id,
			"chat_id":    sender.GetChatID(),
			"bot_id":     bot_id,
			// "uuid":       utils.GenUUID(),
		}
		for k, v := range sender.GetExpandMessageInfo() {
			msg[k] = v
		}
		if sender.f.umod { //订阅号模式
			v, loaded := sender.f.gmsgChan.LoadOrStore(user_id, &GMsgChan{})
			ch := v.(*GMsgChan)
			if !loaded {
				// console.Debug("发送创建：", ch.Chan)
				ch.Chan = make(chan map[string]interface{})
			} else {
				// console.Debug("发送加载：", ch.Chan)
			}
			select {
			case ch.Chan <- msg:
				// console.Debug("管道发送：", msg, ch.Chan)
			case <-time.After(time.Second):
				ch.Lock()
				defer ch.Unlock()
				ch.Msgs = append(ch.Msgs, msg)
				// console.Debug("数组发送：", msg)
			}
			return "", nil
		}
		if sender.f.reply == nil { //未设置回复函数
			var any *common.Function
			var one *common.Function
			for _, function := range Functions {
				if function.Reply != nil && function.Reply.Platform == platform {
					if len(function.Reply.BotsID) == 0 {
						any = function
					} else if Contains(function.Reply.BotsID, sender.f.botid) {
						one = function
					}
				}
			}
			if one == nil && any != nil {
				one = any
			}
			if one != nil {
				message_id := ""
				one.Handle(&Faker{
					Type: "message",
				}, func(vm *goja.Runtime) {
					obj := vm.NewObject()
					for k, v := range msg {
						obj.Set(k, v)
					}
					proxy := vm.NewProxy(obj, &goja.ProxyTrapConfig{
						Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
							message_id = fmt.Sprint(value.Export())
							return true
						},
					})
					vm.Set("msg", proxy)
					vm.Set("message", proxy)
					vm.Set("adapter", sender.f)
				})
				return message_id, nil
			} else { //存储消息
				c := MsgChan{
					Msg:  msg,
					Chan: make(chan string),
				}
				if sender.f.destroid {
					return "", errors.New("adapter destroid")
				}
				sender.f.msgChan <- c
				select {
				case id := <-c.Chan:
					return id, nil
				case <-time.After(time.Second * 5):
					close(c.Chan)
					return "", errors.New("获取消息ID超时")
				}
			}
		} else {
			//todo 阻塞延迟异常
			v := sender.f.reply(msg)
			return v, nil
		}
	}
	return "", nil
}

const (
	MESSAGE_ID = "message_id"
	CONETNT    = "content"
	USER_ID    = "user_id"
	CHAT_ID    = "chat_id"
	IS_ADMIN   = "is_admin"
	// UUID       = "message_id"
)

func (sender *CustomSender) Copy() common.Sender {
	new := reflect.Indirect(reflect.ValueOf(interface{}(sender))).Interface().(CustomSender)
	return &new
}

func (sender *CustomSender) Event() map[string]interface{} {
	e := sender.GetVar("event")
	if e == nil {
		return nil
	}
	return e.(map[string]interface{})
}

func (sender *CustomSender) RecallMessage(ps ...interface{}) {
	recalls := []func(){}
	var timeout int
	for _, p := range ps {
		switch p := p.(type) {
		case int, int64:
			timeout = utils.Int(p)
		case string:
			if p != "" {
				recalls = append(recalls, func() {
					sender.Action(H{
						"type":       "delete_message",
						"message_id": p,
					})
				})
			}
		case []string:
			for _, v := range p {
				if v != "" {
					recalls = append(recalls, func() {
						sender.Action(H{
							"type":       "delete_message",
							"message_id": v,
						})
					})
				}
			}
		case [][]string:
			for _, v := range p {
				for _, v2 := range v {
					recalls = append(recalls, func() {
						sender.Action(H{
							"type":       "delete_message",
							"message_id": v2,
						})
					})
				}
			}
		}
	}
	go func() {
		if timeout != 0 {
			time.Sleep(time.Millisecond * time.Duration(timeout))
		}
		for _, recall := range recalls {
			recall()
		}
	}()
}

func (sender *CustomSender) GroupKick(uid string, reject_add_request bool) error {
	_, err := sender.Reply(mystr.BuildCQCode("kick", H{"user_id": uid, "chat_id": sender.GetChatID(), "forever": reject_add_request}, ""))
	return err
}

func (sender *CustomSender) GroupBan(uid string, duration int) error {
	_, err := sender.Reply(mystr.BuildCQCode("ban", H{"user_id": uid, "chat_id": sender.GetChatID(), "duration": duration}, ""))
	return err
}

func (sender *CustomSender) GroupUnban(uid string) error {
	_, err := sender.Reply(mystr.BuildCQCode("ban", H{"user_id": uid, "chat_id": sender.GetChatID(), "duration": 0}, ""))
	return err
}

func (sender *CustomSender) IsAdmin() bool {
	if sender.f.isAdmin == nil {
		return Contains(strings.Split(MakeBucket(sender.f.botplt).GetString("masters"), "&"), sender.GetUserID())
	}
	return sender.f.isAdmin(sender.GetUserID())
}
