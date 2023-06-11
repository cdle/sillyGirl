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
	Msg  map[string]string
}

type Factory struct {
	botid         string
	botplt        string
	uuid          string
	msgChan       chan MsgChan
	demo          *CustomSender
	reply         func(map[string]string) string
	lm            chan bool
	nm            int64
	recallMessage func(interface{}) bool
	groupKick     func(uid string, gid string, reject_add_request bool) bool
	groupBan      func(uid string, gid string, duration int) bool
	groupUnban    func(uid string, gid string) bool
	isAdmin       func(string) bool
	vm            *goja.Runtime
	ctx           context.Context
	cancel        context.CancelFunc
	destroid      bool
	errorTimes    int
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

func (f *Factory) Init(botplt, botid string) {
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
		go v.Destroy()
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

func (f *Factory) Destroy() {
	BotsLocker.Lock()
	defer BotsLocker.Unlock()
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

func (f *Factory) SetReplyHandler(function func(map[string]string) string) {
	f.reply = func(m map[string]string) string {
		if f.uuid != "" {
			mutex := GetMutex(f.uuid)
			mutex.Lock()
			defer mutex.Unlock()
		}
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").SetReply error:", err)
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

func (f *Factory) GetReplyMessage() *goja.Promise {
	promise, resolve, reject := f.vm.NewPromise()
	go func() {
		select {
		case <-f.ctx.Done():
			logs.Debug("%s adapter %s destroied", f.botplt, f.botid)
			reject("adapter destroied")
		case mc := <-f.msgChan:
			obj := f.vm.NewObject()
			for k, v := range mc.Msg {
				obj.Set(k, v)
			}
			resolve(f.vm.NewProxy(obj, &goja.ProxyTrapConfig{
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
	}()
	return promise
}

func (f *Factory) Send(function func(map[string]string) string) {
	f.SetReplyHandler(function)
}

func (f *Factory) SetRecallMessage(function func(interface{}) bool) {
	f.recallMessage = func(i interface{}) bool {
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").recall error:", err)
			}
		}()
		return function(i)
	}
}

func (f *Factory) SetGroupKick(function func(uid string, gid string, reject_add_request bool) bool) {
	f.groupKick = func(uid string, gid string, reject_add_request bool) bool {
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").GroupKick error:", err)
			}
		}()
		return function(uid, gid, reject_add_request)
	}
}

func (f *Factory) SetGroupBan(function func(uid string, gid string, duration int) bool) {
	f.groupBan = func(uid string, gid string, duration int) bool {
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").SetGroupBan error:", err)
			}
		}()
		return function(uid, gid, duration)
	}
}

func (f *Factory) SetGroupUnban(function func(uid string, gid string) bool) {
	f.groupUnban = func(uid string, gid string) bool {
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").SetgroupUnban error:", err)
			}
		}()
		return function(uid, gid)
	}
}

func (f *Factory) SetIsAdmin(function func(string) bool) {
	f.isAdmin = func(uid string) bool {
		defer func() {
			err := recover()
			if err != nil {
				console.Error("Sender(\""+f.botplt+"\").SetAdmin error:", err)
			}
		}()
		return function(uid)
	}
}

func (f *Factory) Receive(wt interface{}) *CustomSender {
	var demo = *f.demo
	sender := &demo
	props := wt.(map[string]interface{})
	for i := range props {
		switch strings.ToLower(i) {
		case "content":
			sender.details.Content = fmt.Sprint(props[i])
		case "message_id", "messageId":
			sender.details.MessageID = fmt.Sprint(props[i])
		case "user_id", "userId":
			sender.details.UserID = fmt.Sprint(props[i])
		case "chat_id", "chatId", "group_id", "groupId", "group_code", "groupCode":
			sender.details.ChatID = fmt.Sprint(props[i])
		case "user_name", "userName":
			sender.details.Username = fmt.Sprint(props[i])
		case "chat_name", "chatName", "groupName", "group_name":
			sender.details.Chatname = fmt.Sprint(props[i])
		}
	}
	if sender.details.Content != "" {
		Messages <- sender
	}
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

func (sender *CustomSender) Reply(msgs ...interface{}) (string, error) {
	var push = false
	var content = ""
	var platform = sender.f.botplt
	var bot_id = sender.f.botid
	for _, item := range msgs {
		switch item := item.(type) {
		case PUSH:
			push = true
		case string:
			content = item
		}
	}
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
		msg := map[string]string{
			"message_id": sender.GetMessageID(),
			"content":    content,
			"user_id":    sender.GetUserID(),
			"chat_id":    sender.GetChatID(),
			"bot_id":     bot_id,
			// "uuid":       utils.GenUUID(),
		}
		if sender.f.reply == nil {
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
					Type: "*",
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
			} else {
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
					return "", errors.New("get message_id timeout")
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

func (sender *CustomSender) RecallMessage(ps ...interface{}) error {
	if sender.f.recallMessage == nil {
		return nil
	}
	for _, p := range ps {
		switch p := p.(type) {
		case string:
			sender.f.recallMessage(p)
		case []string:
			for _, v := range p {
				sender.f.recallMessage(v)
			}
		case [][]string:
			for _, v := range p {
				for _, v2 := range v {
					sender.f.recallMessage(v2)
				}
			}
		}
	}
	return nil
}

func (sender *CustomSender) GroupKick(uid string, reject_add_request bool) {
	if sender.f.groupKick == nil {
		return
	}
	sender.f.groupKick(uid, sender.GetChatID(), reject_add_request)
}

func (sender *CustomSender) GroupBan(uid string, duration int) {
	if sender.f.groupBan == nil {
		return
	}
	sender.f.groupBan(uid, sender.GetChatID(), duration)
}

func (sender *CustomSender) GroupUnban(uid string) {
	if sender.f.groupUnban == nil {
		return
	}
	sender.f.groupUnban(uid, sender.GetChatID())
}

func (sender *CustomSender) IsAdmin() bool {
	if sender.f.isAdmin == nil {
		return Contains(strings.Split(MakeBucket(sender.f.botplt).GetString("masters"), "&"), sender.GetUserID())
	}
	return sender.f.isAdmin(sender.GetUserID())
}
