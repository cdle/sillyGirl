package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

var cgs []CarryGroup

type CarryGroupsResult struct {
	Success bool         `json:"success"`
	Data    []CarryGroup `json:"data"`
	Page    int          `json:"page"`
	Total   int          `json:"total"`
	Time    time.Time    `json:"time"`
}

var CarryGroups = MakeBucket("CarryGroups")

var carryCounter int64

type QMessage struct {
	UserID    string        `json:"user_id"`
	Content   string        `json:"content"`
	MessageID string        `json:"message_id"`
	From      common.Sender `json:"-"`
	To        *Factory      `json:"-"`
}

// LOGIC
func initCarry() {
	AddCommand([]*common.Function{
		{
			Rules:    []string{`raw [\s\S]*`},
			Hidden:   true,
			Priority: 9999,
			Handle: func(s common.Sender, _ func(vm *goja.Runtime)) interface{} {
				var bot_id = s.GetBotID()
				var platform = s.GetImType()
				var chat_id = s.GetChatID()
				var user_id = s.GetUserID()
				var content = s.GetContent()
				var message_id = s.GetMessageID()
				var from *CarryGroup //判断当前消息来自采集源
				var cgs = cgs
				var uuid = fmt.Sprintf("%d. ", atomic.AddInt64(&carryCounter, 1))
				var ss = &Strings{}
				var event = s.Event()
				if event != nil {
					if event["type"] == "delete_message" {
						queues.Range(func(key, value any) bool {
							q := value.(*Queue)
							for _, qm := range q.GetValues() {
								if qm.From != nil && qm.From.GetMessageID() == event["message_id"] {
									qm.To.Sender().RecallMessage(qm.MessageID)
								}
							}
							return true
						})
					}
					s.Continue()
					return nil
				}
				for i := range cgs {
					if chat_id == cgs[i].ID && cgs[i].In && cgs[i].Enable {
						from = &cgs[i]
						break
					}
				}
				if nil == from { //非采集群
					s.Continue()
					return nil
				}
				//采集消息去重逻辑
				q := NewQueue(chat_id, 50)
				if from.Deduplication {
					for _, qm := range q.GetValues() {
						v := ss.HansSimilarity(qm.Content, content)
						if v > 0.9 {
							console.Debug("%s 忽略重复采集信息", uuid)
							return nil
						}
					}
				}
				_ = q.Enqueue(&QMessage{
					UserID:    user_id,
					Content:   content,
					MessageID: message_id,
				})
				bots_id := GetAdapterBotsID(platform)
				if len(from.BotsID) == 0 && len(bots_id) != 0 {
					from.BotsID = bots_id
				}
				//判断是否来自指定采集机器人
				var from_right_bot bool
				if len(from.BotsID) != 0 && from.BotsID[0] == bot_id {
					from_right_bot = true
				}
				//检测指定机器人是否离线，离线则使用其他第一个机器人，否则忽略消息
				if !from_right_bot {
					if len(from.BotsID) != 0 {
						if Contains(bots_id, bot_id) {
							console.Debug("%s 忽略机器人(%s)消息非采集指定机器人(%s)消息", uuid, bot_id, from.BotsID[0])
							return nil
						}
					}
					if len(bots_id) != 0 && bots_id[0] != bot_id { //不是第一个机器人
						console.Debug("%s 忽略机器人(%s)消息非其他第一个机器人(%s)的消息", uuid, bot_id, bots_id[0])
						return nil
					}
				}
				console.Debug("%s 当前采集群 %s", uuid, chat_id)
				//预测采集白名单、黑名单
				if len(from.Allowed) != 0 { //白名单
					if !Contains(from.Allowed, user_id) {
						console.Debug("%s 用户(%s)不在采集群白名单 %v", uuid, chat_id)
						return nil
					}
				} else {
					if Contains(from.Prohibited, user_id) {
						console.Debug("%s 用户(%s)在采集群黑名单 %v", uuid, chat_id)
						return nil
					}
				}
				//预测采集包含、排除词
				if len(from.Include) != 0 { //包含
					if word := Include(content, from.Include); word == "" {
						console.Debug("%s 消息中无采集包含词", uuid)
						return nil
					}
				}
				if len(from.Exclude) != 0 { //排除
					if word := Include(content, from.Exclude); word != "" {
						console.Debug("%s 消息中有采集排除词 %s", uuid, word)
						return nil
					}
				}
				var outs []CarryGroup //预测转发群
				for i := range cgs {
					if cgs[i].Enable && cgs[i].Out && cgs[i].ID != chat_id {
						for j := range cgs[i].From {
							if cgs[i].From[j] == chat_id {
								if len(cgs[i].Allowed) != 0 { //白名单
									if !Contains(cgs[i].Allowed, user_id) {
										console.Debug("%s 用户(%s)不在转发群(%s)白名单 %v", uuid, user_id, cgs[i].ID)
										continue
									}
								} else {
									if Contains(cgs[i].Prohibited, user_id) {
										console.Debug("%s 用户(%s)在转发群(%s)黑名单 %v", uuid, user_id, cgs[i].ID)
										continue
									}
								}
								if len(cgs[i].Include) != 0 { //包含
									if word := Include(content, cgs[i].Include); word == "" {
										console.Debug("%s 消息中无转发(s)包含词", uuid, cgs[i].ID)
										continue
									}
								}
								if len(cgs[i].Exclude) != 0 { //排除
									if word := Include(content, cgs[i].Exclude); word != "" {
										console.Debug("%s 消息中有转发(s)排除词 %s", uuid, cgs[i].ID, word)
										continue
									}
								}
								outs = append(outs, cgs[i])
							}
						}
					}
				}
				num := len(outs)
				console.Debug("%s 预测转发群数目 %v", uuid, num)
				if num == 0 {
					return nil
				}
				var scripts = []string{}
				//执行采集脚本
				fs := Functions
				for j := range from.Scripts {
					for i := range fs {
						if fs[i].UUID == from.Scripts[j] && !Contains(scripts, fs[i].UUID) {
							fs[i].Handle(s, nil)
							content = s.GetContent()
							if content == "" {
								return nil
							}
							scripts = append(scripts, fs[i].UUID)
						}
					}
				}
				//执行转发脚本
				for i := range outs {
					var scripts = scripts
					var content = content
					for j := range outs[i].Scripts {
						for k := range fs {
							if fs[k].UUID == outs[i].Scripts[j] && !Contains(scripts, fs[k].UUID) { //
								fs[k].Handle(s, nil)
								content = s.GetContent()
								// fmt.Println(content)
								if content == "" {
									goto HELL
								}
								scripts = append(scripts, fs[k].UUID)
							}
						}
					}
				HELL:
					if content != "" { //选择机器人
						platform := outs[i].Platform
						chat_id := outs[i].ID
						adapter, err := GetAdapter(platform, outs[i].BotsID...)
						if adapter == nil {
							console.Warn("%s (%s)转发群(%s)相关机器人%v都不在线", uuid, platform, chat_id, outs[i].BotsID)
							continue
						}
						if err != nil {
							console.Debug("%s 指定(%s)机器人都不在线，转发群(%s)已选择其他机器人(%s)推送", uuid, platform, chat_id, adapter.botid)
						}
						if adapter != nil {
							//采集消息去重逻辑
							q := NewQueue(chat_id, 50)
							if outs[i].Deduplication {
								for _, qm := range q.GetValues() {
									v := ss.HansSimilarity(qm.Content, content)
									if v > 0.9 {
										console.Debug("%s 忽略重复转发信息", uuid)
										continue
									}
								}
							}
							qm := &QMessage{
								UserID:  user_id,
								Content: content,
								From:    s,
								To:      adapter,
							}
							_ = q.Enqueue(qm)
							message_id, err := adapter.Push(map[string]string{
								CONETNT: content,
								CHAT_ID: chat_id,
							})
							if err == nil {
								qm.MessageID = message_id
							}
						}
					}
				}
				return nil
			},
		},
	})

	setCgs()
	storage.Watch(CarryGroups, nil, func(old, new, key string) *storage.Final {
		console.Log("已更新搬运数据")
		ocg := CarryGroup{}
		ncg := CarryGroup{}
		json.Unmarshal([]byte(old), &ocg)
		json.Unmarshal([]byte(new), &ncg)
		tmp := cgs
		if old != "" {
			if new == "" { // 删除
				if ocg.ID != "" {
					for i, cg := range tmp {
						if cg.ID == ocg.ID {
							tmp = append(tmp[:i], tmp[i+1:]...)
							name := cg.ChatName
							if name == "" {
								name = cg.ID
							}
							RemListenOnGroup(cg.ID, fmt.Sprintf("已为采集群(%s)关闭监听模式", name))
							break
						}
					}
				} else {
					return nil
				}
			} else { // 修改
				if ocg.ID != "" {
					for i, cg := range tmp {
						if cg.ID == ocg.ID {
							tmp[i] = ncg
							name := ncg.ChatName
							if name == "" {
								name = ncg.ID
							}
							if ncg.In {
								if ncg.Enable {
									AddListenOnGroup(ncg.ID, fmt.Sprintf("已为采集群(%s)开启监听模式", name), ncg.Platform)
									AddNoReplyGroups(ncg.ID, fmt.Sprintf("已为采集群(%s)开启禁言模式", name), ncg.Platform)
								} else {
									RemListenOnGroup(ncg.ID, fmt.Sprintf("已为采集群(%s)关闭监听模式", name))
								}
							} else {
								RemListenOnGroup(ncg.ID, fmt.Sprintf("已为采集群(%s)关闭监听模式", name))
							}
							break
						}
					}
				} else {
					return nil
				}
			}
		} else { //创建
			if ncg.ID != "" {
				tmp = append(tmp, ncg)
				if ncg.In && ncg.Enable {
					name := ncg.ChatName
					if name == "" {
						name = ncg.ID
					}
					AddListenOnGroup(ncg.ID, fmt.Sprintf("已为采集群(%s)开启监听模式", name), ncg.Platform)
					AddNoReplyGroups(ncg.ID, fmt.Sprintf("已为采集群(%s)开启禁言模式", name), ncg.Platform)
				}
			} else {
				return nil
			}
		}
		sort.Sort(byCreatedAt(tmp))
		for i := range tmp {
			tmp[i].Index = i + 1
		}
		cgs = tmp
		return nil
	})
}

func setCgs() {
	CarryGroups.Foreach(func(b1, b2 []byte) error {
		cg := CarryGroup{}
		err := json.Unmarshal(b2, &cg)
		if err != nil {
			return nil
		}
		if cg.In && cg.Enable {
			name := cg.ChatName
			if name == "" {
				name = cg.ID
			}
			AddListenOnGroup(cg.ID, fmt.Sprintf("已为采集群(%s)开启监听模式", name), cg.Platform)
			AddNoReplyGroups(cg.ID, fmt.Sprintf("已为采集群(%s)开启禁言模式", name), cg.Platform)
		}
		cgs = append(cgs, cg)
		return nil
	})
	sort.Sort(byCreatedAt(cgs))
	for i := range cgs {
		cgs[i].Index = i + 1
	}
}

type CarryGroup struct {
	Index          int      `json:"id"`             //编号 顺序编号
	In             bool     `json:"in"`             //搬进来 勾选按钮
	Out            bool     `json:"out"`            //运出去 勾选按钮
	From           []string `json:"from"`           //采集源
	Allowed        []string `json:"allowed"`        //白名单模式
	Prohibited     []string `json:"prohibited"`     //黑名单模式 Select选择器多选
	ID             string   `json:"chat_id"`        //群组ID 文字表单
	ChatName       string   `json:"chat_name"`      //群昵称 文字表单
	Remark         string   `json:"remark"`         //备注
	Platform       string   `json:"platform"`       //平台 Select选择器单选
	Enable         bool     `json:"enable"`         //启用状态 开关
	Include        []string `json:"include"`        //包含关键词 多个关键词用逗号隔开 用户复制粘贴过去后自动转换成多彩标签
	Exclude        []string `json:"exclude"`        //排除关键词 包含关键词
	CreatedAt      int      `json:"created_at"`     //创建时间戳(秒)转换成日期
	BotsID         []string `json:"bots_id"`        //工作机器人 多选
	Scripts        []string `json:"scripts"`        //处理脚本
	Deduplication  bool     `json:"deduplication"`  //文本去重
	Deduplication2 bool     `json:"deduplication2"` //图片去重
}

// CARRY API
func init() {
	GinApi(GET, "/api/carry/groups", RequireAuth, func(ctx *gin.Context) {
		current := utils.Int(ctx.Query("current"))
		pageSize := utils.Int(ctx.Query("pageSize"))
		rr := CarryGroupsResult{
			Success: true,
		}
		cgs := cgs
		rr.Total = len(cgs)
		if current == 0 {
			current = 1
		}
		if pageSize == 0 {
			pageSize = 20
		}
		begin := (current - 1) * pageSize
		end := (current) * pageSize
		if end > rr.Total {
			end = rr.Total
		}
		if begin > end {
			begin = end
		}
		rr.Data = cgs[begin:end]
		for i := range rr.Data {
			gn := &Nickname{
				ID: rr.Data[i].ID,
			}
			nickname.First(gn)
			if gn.Value != "" {
				rr.Data[i].ChatName = gn.Value
			}
		}
		ctx.JSON(200, rr)
	})
	GinApi(GET, "/api/carry/group_names", RequireAuth, func(ctx *gin.Context) {
		cgs := cgs
		var names = map[string]string{}
		for _, cg := range cgs {
			names[cg.ID] = cg.ChatName
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    names,
		})
	})
	GinApi(GET, "/api/proxy/scripts", RequireAuth, func(ctx *gin.Context) {
		var scripts = map[string]string{}
		functions := Functions
		for _, function := range functions {
			if function.UUID != "" {
				scripts[function.UUID] = function.Title + ".js"
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    scripts,
		})
	})
	var isNumeric = func(keyword string) bool {
		for _, c := range keyword {
			if c != '.' && (c < '0' || c > '9') {
				return false
			}
		}
		return true
	}
	GinApi(GET, "/api/proxy/rules", RequireAuth, func(ctx *gin.Context) {
		keyword := ctx.Query("keyword")
		var scripts = map[string]string{}
		scripts[keyword] = keyword
		if strings.HasSuffix(keyword, ".") && !isNumeric(keyword) {
			for _, suffix := range []string{"com", "cn"} {
				scripts[keyword+suffix] = keyword + suffix
			}
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data":    scripts,
		})
	})
	GinApi(GET, "/api/carry/group_selects", RequireAuth, func(ctx *gin.Context) {
		chat_id := ctx.Query("chat_id")
		platform := ctx.Query("platform")
		cgs := cgs
		var names = map[string]string{}
		var bots_id = []string{}
		var users = []string{}
		for _, cg := range cgs {
			if cg.In {
				if cg.ChatName != "" {
					names[cg.ID] = cg.ChatName
				} else {
					if cg.Remark != "" {
						names[cg.ID] = cg.Remark
					} else {
						names[cg.ID] = cg.ID
					}
				}
			}
			if cg.ID == chat_id {
				users = append(users, cg.Allowed...)
				users = append(users, cg.Prohibited...)
				if platform == "" {
					platform = cg.Platform
				}
			}
		}
		bots_id = GetAdapterBotsID(platform)
		var scripts = map[string]string{}
		functions := Functions
		for _, function := range functions {
			if function.UUID != "" && ((len(function.Rules) == 0 && !function.OnStart && !function.Module && len(function.Https) == 0 && function.Reply == nil) || function.Carry) {
				scripts[function.UUID] = function.Title + ".js"
			}
		}
		var user_names = []NicklabeL{}
		nickname.Foreach(func(b1, b2 []byte) error {
			v := &Nickname{}
			code := string(b1)
			err := json.Unmarshal(b2, v)
			if err == nil {
				platforms = append(platforms, v.Platform)
				if Contains(users, code) {
					user_names = append(user_names, NicklabeL{
						Label: fmt.Sprintf("%s(%s)", v.Value, code),
						Value: code,
					})
				}
			}
			return nil
		})
		ctx.JSON(200, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"user_names":  user_names,
				"group_names": names,
				"bots_id":     bots_id,
				"platforms":   getPltsArray(),
				"scripts":     scripts,
			},
		})

	})
	GinApi(POST, "/api/carry/group", RequireAuth, func(ctx *gin.Context) {
		// 将请求的 JSON 数据解析为一个 map[string]interface{} 类型的变量
		var updateData map[string]interface{}
		err := ctx.BindJSON(&updateData)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		v, ok := updateData["chat_id"]
		if !ok {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "群号不能为空",
			})
			return
		}
		chat_id := v.(string)
		var cg = CarryGroup{
			ID: chat_id,
		}
		CarryGroups.First(&cg)
		// if err != nil {
		// 	ctx.JSON(200, map[string]interface{}{
		// 		"success":      false,
		// 		"errorMessage": err.Error(),
		// 	})
		// 	return
		// }
		for key, value := range updateData {
			switch key {
			case "in":
				if in, ok := value.(bool); ok {
					cg.In = in
				}
			case "out":
				if out, ok := value.(bool); ok {
					cg.Out = out
				}
			case "deduplication":
				if deduplication, ok := value.(bool); ok {
					cg.Deduplication = deduplication
				}
			case "deduplication2":
				if deduplication, ok := value.(bool); ok {
					cg.Deduplication2 = deduplication
				}
			case "from":
				if from, ok := value.([]interface{}); ok {
					cg.From = toStringSlice(from)
				}
			case "allowed":
				if allowed, ok := value.([]interface{}); ok {
					cg.Allowed = toStringSlice(allowed)
				}
			case "prohibited":
				if prohibited, ok := value.([]interface{}); ok {
					cg.Prohibited = toStringSlice(prohibited)
				}
			case "chat_name":
				if chatName, ok := value.(string); ok {
					cg.ChatName = chatName
				}
			case "remark":
				if remark, ok := value.(string); ok {
					cg.Remark = remark
				}
			case "platform":
				if platform, ok := value.(string); ok {
					cg.Platform = platform
				}
			case "enable":
				if disable, ok := value.(bool); ok {
					cg.Enable = disable
				}
			case "include":
				if include, ok := value.([]interface{}); ok {
					cg.Include = toStringSlice(include)
				}
			case "exclude":
				if exclude, ok := value.([]interface{}); ok {
					cg.Exclude = toStringSlice(exclude)
				}
			case "bots_id":
				if botsID, ok := value.([]interface{}); ok {
					cg.BotsID = toStringSlice(botsID)
				}
			case "scripts":
				if scripts, ok := value.([]interface{}); ok {
					cg.Scripts = toStringSlice(scripts)
				}
			}
		}
		if cg.CreatedAt == 0 {
			cg.CreatedAt = int(time.Now().Unix())
		}
		CarryGroups.Set(chat_id, utils.JsonMarshal(cg))
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
	GinApi(DELETE, "/api/carry/group", RequireAuth, func(ctx *gin.Context) {
		cg := &CarryGroup{}
		err := ctx.BindJSON(cg)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": err.Error(),
			})
			return
		}
		if cg.ID == "" {
			ctx.JSON(200, map[string]interface{}{
				"success":      false,
				"errorMessage": "群号不为空",
			})
			return
		}
		CarryGroups.Set(cg.ID, "")
		ctx.JSON(200, map[string]interface{}{
			"success": true,
		})
	})
}

type byCreatedAt []CarryGroup

func (s byCreatedAt) Len() int {
	return len(s)
}

func (s byCreatedAt) Less(i, j int) bool {
	return s[i].CreatedAt > s[j].CreatedAt
}

func (s byCreatedAt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 将 []interface{} 转为 []string 的工具函数
func toStringSlice(intfSlice []interface{}) []string {
	stringSlice := make([]string, len(intfSlice))
	for i, intf := range intfSlice {
		if str, ok := intf.(string); ok {
			stringSlice[i] = str
		}
	}
	return stringSlice
}

func Contains(strs []string, str ...string) bool {
	for _, s := range str {
		for _, str := range strs {
			if s == str {
				return true
			}
		}
	}
	return false
}

func Include(content string, includes []string) string {
	for _, include := range includes {
		if len(include) > 2 && include[0] == '/' && include[len(include)-1] == '/' {
			pattern := include[1 : len(include)-1]
			_, err := regexp.Compile(pattern)
			if err != nil {
				console.Error("包含词/排除词正则表达式 %s 错误 %s", include, err.Error())
				continue
			}
			match, err := regexp.MatchString(pattern, content)
			if err == nil && match {
				return include
			}
		} else {
			if strings.Contains(content, include) {
				return include
			}
		}
	}
	return ""
}
