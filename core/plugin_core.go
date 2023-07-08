package core

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/common"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/robfig/cron/v3"
)

type myFieldNameMapper struct{}

var mutexMap = make(map[string]*sync.Mutex)
var mutexMapMutex sync.Mutex

func GetMutex(uuid string) *sync.Mutex {
	mutexMapMutex.Lock()
	defer mutexMapMutex.Unlock()

	if mutex, ok := mutexMap[uuid]; ok {
		return mutex
	}

	mutex := &sync.Mutex{}
	mutexMap[uuid] = mutex
	return mutex
}

func (tfm myFieldNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(`json`)
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if parser.IsIdentifier(tag) {
		return tag
	}
	return uncapitalize(f.Name)
}

func uncapitalize(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func (tfm myFieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return uncapitalize(m.Name)
}

var RegistFuncs = map[string]interface{}{}

var plugins = MakeBucket("plugins")

type Route struct {
	Path      string  `json:"path"`
	Name      string  `json:"name"`
	Component string  `json:"component,omitempty"`
	Routes    []Route `json:"routes,omitempty"`
	// Key       string  `json:"key,omitempty"`
	CreateAt string `json:"create_at"`
}

func CancelPluginlistening(uuid string) {
	// logs.Info(`k, c.Function, c.Function.Rules`)
	for _, wait := range waits {
		wait.Foreach(func(k int64, c *Carry) bool {
			if uuid == c.UUID {
				c.Chan <- errors.New("uinstall")
			}
			return true
		})
	}
}

var debug = sillyGirl.GetBool("debug", false)

func initPlugins() {
	storage.Watch(sillyGirl, "debug", func(old, new, key string) *storage.Final {
		debug = new == "true"
		return nil
	})
	plugins.Foreach(func(key, data []byte) error {
		f, cbs, err := initPlugin(string(data), string(key))
		if err != nil {
			console.Error("初始化插件%s错误: %v", key, err)
		}
		for _, cb := range cbs {
			cb()
		}
		AddCommand([]*common.Function{f})
		// os.WriteFile(fmt.Sprintf("%s/%s.js", pluginPath, f.Title), data, 0600)
		return nil
	})
	var pluginLock = new(sync.Mutex)
	storage.Watch(plugins, nil, func(old, new, key string) (fin *storage.Final) {
		pluginLock.Lock()
		defer pluginLock.Unlock()
		if new == "install" {
			for _, p := range plugin_list {
				if p.UUID != key {
					continue
				}
				script := string(fetchScript(p.Address, key))
				if f, _, _ := initPlugin(script, p.UUID); f.CreateAt != "" {
					fin = &storage.Final{
						Now: script,
					}
					fin = &storage.Final{}
					su := &ScriptUtils{script: script}
					su.SetValue("origin", p.Organization)
					new = su.script
					break
				} else {
					return &storage.Final{
						Error: errors.New("订阅源异常"),
					}
				}
			}
		}

		if new != "" {
			if new == "reload" {
				new = old
			}
			fin = &storage.Final{}
			su := &ScriptUtils{script: new}
			if version := su.GetValue("version"); version == "" || regexp.MustCompile(`v\d+\.\d+\.\d`).FindString(version) != version {
				su.SetValue("version", "v1.0.0")
			}
			if auhtor := su.GetValue("author"); auhtor == "" {
				su.SetValue("author", "佚名")
			}
			if su.GetValue("description") == "" {
				su.SetValue("description", "🐒这个人很懒什么都没有留下")
			} //module
			title := su.GetValue("title")
			if title == "" {
				su.SetValue("title", "无名脚本")
			}
			if message := GetMessageByUUID(key); message != "" {
				su.SetValue("message", message)
			}
			if title != "无名脚本" && title != "" {
				onStart := su.GetValue("on_start")
				if onStart != "true" {
					module := su.GetValue("module")
					if module != "true" {
						if module != "" {
							su.DeleteValue("module")
						}
						web := su.GetValue("web")
						if web != "true" {
							if web != "" {
								su.DeleteValue("web")
							}
							// rule := su.GetValue("rule")
							// if rule == "" {
							// su.SetValue("rule", title)
							// }
						} else {
							su.DeleteValue("rule")
							su.DeleteValue("cron")
							su.DeleteValue("admin")
							su.DeleteValue("priority")
							su.DeleteValue("platform")
						}
					} else {
						su.DeleteValue("rule")
						su.DeleteValue("cron")
						su.DeleteValue("web")
						su.DeleteValue("admin")
						su.DeleteValue("priority")
						su.DeleteValue("platform")
					}
				} else {
					// su.DeleteValue("rule")
					// su.DeleteValue("cron")
					su.DeleteValue("web")
					su.DeleteValue("admin")
					su.DeleteValue("priority")
					su.DeleteValue("platform")
					su.DeleteValue("module")
				}
			}
			create_at := su.GetValue("create_at")
			if _, err := time.Parse("2006-01-02 15:04:05", create_at); err != nil {
				su.SetValue("create_at", time.Now().Format("2006-01-02 15:04:05"))
			}
			fin.Now = su.script
			if su.script != new {
				fin.Message = su.script
			} else if title != (&ScriptUtils{script: old}).GetValue("title") {
				fin.Message = "标题变更。"
			}
			new = su.script
		}
		f, cbs, err := initPlugin(new, key)
		if err != nil && new != "" {
			pluginConsole(key).Error(err)
		}
		apd := false
		for i := range Functions {
			if Functions[i].UUID == key {
				DestroyAdapterByUUID(key)
				Functions[i].Running = false
				if len(Functions[i].CronIds) != 0 {
					for _, id := range Functions[i].CronIds {
						CRON.Remove(cron.EntryID(id))
					}
				}
				Functions = append(Functions[:i], Functions[i+1:]...)
				CancelPluginCrons(key)
				CancelPluginWebs(key)
				CancelPluginlistening(key)
				CancelHttpListen(key)
				remStatic(key)
				storage.DisableHandle(key)
				if new != "" {
					AddCommand([]*common.Function{f})
					if old == "" {
						console.Log("已加载 %s.js", f.Title)
					} else if !f.OnStart {
						console.Log("已重载 %s.js", f.Title)
					}
				} else {
					of, _, _ := initPlugin(old, key)
					console.Log("已卸载 %s.js", of.Title)
				}
				apd = true
				break
			}
		}
		for _, cb := range cbs {
			cb()
		}
		if !apd {
			AddCommand([]*common.Function{f})
		}
		if f.UUID != "" && f.Public {
			go func() {
				os.WriteFile(fmt.Sprintf("%s/%s.js", plugin_download_file, f.UUID), []byte(publicScript(plugins.GetString(f.UUID))), 0666)
				os.WriteFile(plugin_path+"list.json", utils.JsonMarshal(GetPublicResponse()), 0666)
			}()
		}
		return
	})
}

func initPlugin(data string, uuid string) (*common.Function, []func(), error) {
	var cbs = []func(){}
	var err error
	var rules []string
	var imType *common.Filter
	var userId *common.Filter
	var groupId *common.Filter
	var cron = map[string]string{}
	var admin bool
	var disable bool
	var priority int
	var title string
	var public bool
	var description string
	var icon string
	var version string = "v1.0.0"
	var author string
	var create_at string
	var module bool
	var web bool
	var encrypt bool
	var onStart bool
	var origin = "自定义"
	var https = []*common.Http{}
	var message *common.Reply
	var FindAll bool
	var hasForm bool
	var carry bool
	var classes = []string{}
	ks := map[string]bool{}
	ress := regexp.MustCompile(
		`\*\s?@([\d\w+-]+)\s+([^\n]+?)\n`,
	).FindAllStringSubmatch(data, -1)
	for _, res := range ress {
		switch res[1] {
		case "rule", "match", "regex", "pattern":
			rule := strings.TrimSpace(res[2])
			rule = parseReply3(rule, func(s1, s2 string) {
				k := s1 + "." + s2
				if _, ok := ks[k]; !ok {
					cbs = append(cbs, func() {
						storage.Watch(MakeBucket(s1), s2, func(old, new, key string) *storage.Final {
							return &storage.Final{
								EndFunc: func() {
									plugins.Set(uuid, "reload")
								},
							}
						}, uuid)
					})
					ks[k] = true
				}
			})
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
			if len(_rs) != 0 {
				rules = append(rules, _rs...)
			} else {
				rules = append(rules, rule)
			}
		case "class":
			classes = append(classes, regexp.MustCompile(`[\S]+`).FindAllString(res[2], -1)...)
			classes = utils.Unique(classes)
		case "platform", "imType", "platform+", "imType+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			imType = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "platform-", "imType-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			imType = &common.Filter{
				BlackMode: true,
				Items:     item,
			}
		case "userId", "userID", "uid", "userId+", "userID+", "uid+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			userId = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "userId-", "userID-", "uid-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			userId = &common.Filter{
				BlackMode: true,
				Items:     item,
			}
		case "groupId", "groupID", "groupCode", "chat_id", "chat_id+", "chatId", "chatID", "gid", "groupId+", "groupID+", "groupCode+", "chatId+", "chatID+", "gid+":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			groupId = &common.Filter{
				BlackMode: false,
				Items:     item,
			}
		case "groupId-", "groupID-", "groupCode-", "chatId-", "chat_id-", "chatID-", "gid-":
			var item []string
			for _, i := range regexp.MustCompile(`[\d\w-]+`).FindAllString(res[2], -1) {
				item = append(item, strings.TrimSpace(i))
			}
			groupId = &common.Filter{
				BlackMode: true,
				Items:     item,
			}

		case "admin":
			admin = strings.TrimSpace(res[2]) == "true"
		case "disable":
			disable = strings.TrimSpace(res[2]) == "true"
		case "findall":
			FindAll = strings.TrimSpace(res[2]) == "true"
		case "priority":
			priority = utils.Int(strings.TrimSpace(res[2]))
		case "title", "name", "show":
			title = strings.TrimSpace(res[2])
		case "public":
			public = strings.TrimSpace(res[2]) == "true"
		case "description":
			description = strings.TrimSpace(res[2])
		case "icon":
			icon = strings.TrimSpace(res[2])
		case "version":
			version = strings.TrimSpace(res[2])
		case "author":
			author = strings.TrimSpace(res[2])
		case "http":
			ss := regexp.MustCompile(`[\S]+`).FindAllString(strings.TrimSpace(res[2]), -1)
			if len(ss) == 2 {
				https = append(https, &common.Http{
					Path:   ss[1],
					Method: strings.ToUpper(ss[0]),
				})
			} else {
				console.Warn("http param is not 2")
			}
		case "message":
			ss := regexp.MustCompile(`[\S]+`).FindAllString(strings.TrimSpace(res[2]), -1)
			if len(ss) > 1 {
				if len(ss) == 2 && ss[1] == "*" {
					message = &common.Reply{
						Platform: ss[0],
						BotsID:   []string{},
					}
				} else {
					message = &common.Reply{
						Platform: ss[0],
						BotsID:   ss[1:],
					}
				}

			} else {
				console.Warn("message param is 0")
			}
		case "create_at":
			create_at = strings.TrimSpace(res[2])
		case "origin":
			origin = strings.TrimSpace(res[2])
		case "module":
			module = strings.TrimSpace(res[2]) == "true"
		case "web":
			web = strings.TrimSpace(res[2]) == "true"
		case "carry":
			carry = strings.TrimSpace(res[2]) == "true"
		case "encrypt":
			encrypt = strings.TrimSpace(res[2]) == "true"
		case "on_start":
			onStart = strings.TrimSpace(res[2]) == "true"
		case "form":
			hasForm = true
		case "paterner":
			paterner := strings.TrimSpace(res[2])
			go func() {
				time.Sleep(time.Second * 2)
				getPaterner(uuid, strings.TrimSpace(paterner))
			}()
		default:
			cron_ := strings.TrimSpace(res[2])
			cron_ = strings.ReplaceAll(cron_, `\/`, "/")
			if strings.HasPrefix(res[1], "cron") {
				cron[res[1]] = cron_
			}
		}
	}
	script := ""
	if encrypt {
		script = DecryptPlugin(string(data))
	} else {
		script = string(data)
	}
	script = halfDeEct(script)
	script = strings.ReplaceAll(script, "new Sender", "Sender")
	script = strings.ReplaceAll(script, "new Bucket", "Bucket")
	prg, err2 := goja.Compile(title+".js", script, false)
	if err == nil && err2 != nil {
		err = err2
	}
	// if err == nil && len(rules) == 0 && cron != "" {
	// 	err = fmt.Errorf("无效的脚本%s", title)
	// }
	if web {
		onStart = true
	}
	// if icon == "" {
	// 	icon = "https://joeschmoe.io/api/v1/random?t=" + fmt.Sprint(time.Now().Nanosecond())
	// }
	var running func() bool

	f := &common.Function{
		Handle: func(s common.Sender, set func(vm *goja.Runtime)) interface{} {
			if !debug {
				defer func() {
					err := recover()
					if err != nil {
						pluginConsole(uuid).Error(err)
						// s.Reply(fmt.Sprint(err))
					}
				}()
			}
			if err2 != nil {
				panic(err2)
			}
			loop := eventloop.NewEventLoop()
			loop.Run(func(vm *goja.Runtime) {
				SetPluginMethod(vm, uuid, onStart, running)
				ss := &SenderJsIplm{
					Message:    s,
					Vm:         vm,
					Private:    "private",
					Group:      "group",
					Routine:    "routine",
					Persistent: "persistent",
					UUID:       uuid,
				}
				vm.Set("msg", goja.Undefined())
				vm.Set("message", goja.Undefined())
				vm.Set("res", goja.Undefined())
				vm.Set("req", goja.Undefined())
				vm.Set("sender", ss)
				vm.Set("run", func(uuid string) bool { //执行子脚本
					fs := Functions
					for i := range fs {
						if fs[i].UUID == uuid {
							fs[i].Handle(s, nil)
							return true
						}
					}
					return false
				})
				vm.Set("s", ss)
				vm.Set("InitAdapter", func(plt, botid string, params map[string]interface{}) *Factory {
					f := &Factory{
						uuid: uuid,
						vm:   vm,
					}
					f.Init(plt, botid, params)
					return f
				})
				vm.Set("initAdapter", func(plt, botid string, params map[string]interface{}) *Factory {
					f := &Factory{
						uuid: uuid,
						vm:   vm,
					}
					f.Init(plt, botid, params)
					return f
				})
				getAdapter := func(plt, botid string) map[string]interface{} {
					adapter, err := GetAdapter(plt, botid)
					errstr := ""
					if err != nil {
						errstr = err.Error()
					}
					return map[string]interface{}{
						"error":   errstr,
						"adapter": adapter,
					}
				}
				vm.Set("GetAdapter", getAdapter)
				vm.Set("getAdapter", getAdapter)
				vm.Set("getAdapterBotsID", GetAdapterBotsID)
				vm.Set("getAdapterBotPlts", GetAdapterBotPlts)
				vm.Set("GetAdapterBotsID", GetAdapterBotsID)
				vm.Set("GetAdapterBotPlts", GetAdapterBotPlts)
				vm.Set("running", running)
				vm.Set("Running", running)
				vm.Set("uuid", func() string {
					return uuid
				})
				vm.Set("genUuid", func() string {
					return utils.GenUUID()
				})
				vm.Set("genUUID", func() string {
					return utils.GenUUID()
				})
				vm.Set("UUID", func() string {
					return uuid
				})
				if set != nil {
					set(vm)
				}
				_, err := vm.RunProgram(prg)
				if err != nil {
					pluginConsole(uuid).Error(strings.ReplaceAll(strings.ReplaceAll(err.Error(), "node_modules/", ""), "github.com/dop251/goja_nodejs/require", ""))
				}
			})
			return nil
		},
		Rules:       rules,
		ImType:      imType,
		UserId:      userId,
		GroupId:     groupId,
		Cron:        cron,
		Admin:       admin,
		Priority:    priority,
		Disable:     disable,
		UUID:        uuid,
		Title:       title,
		Public:      public,
		Description: description,
		Icon:        icon,
		Version:     version,
		Author:      author,
		CreateAt:    create_at,
		Module:      module,
		Encrypt:     encrypt,
		OnStart:     onStart,
		Origin:      origin,
		Running:     onStart,
		Reply:       message,
		Https:       https,
		FindAll:     FindAll,
		HasForm:     hasForm,
		Carry:       carry,
		Classes:     classes,
	}
	running = func() bool {
		return f.Running
	}

	return f, cbs, err
}

func GetFunctionByUUID(uuid string) *common.Function {
	for _, f := range Functions {
		if f.UUID == uuid {
			return f
		}
	}
	return nil
}

func ChatID(p interface{}) string {
	switch p := p.(type) {
	case int:
		if p == 0 {
			return ""
		} else {
			return utils.Itoa(p)
		}
	case string:
		return p
	case nil:
		return ""
	default:
		return utils.Itoa(p)
	}
}
