package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

var pluginLock = new(sync.Mutex)

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
		pluginLock.Lock()
		defer pluginLock.Unlock()
		f, cbs, err := initPlugin(string(data), string(key), "")
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

	storage.Watch(plugins, nil, func(old, new, key string) (fin *storage.Final) {
		pluginLock.Lock()
		defer pluginLock.Unlock()
		if new == "install" {
			for _, p := range plugin_list {
				if p.UUID != key {
					continue
				}
				if p.Type != "goja" { //下载目录插件
					// Content-Type
					var prefix = "?uuid=" + p.UUID
					address := p.Address
					if !strings.HasSuffix(address, "list.json") {
						address = address + "/api/plugins/download" + prefix
					} else {
						address = strings.ReplaceAll(address, "list.json", "download"+prefix)
					}
					resp, err := http.Get(address)
					if err != nil {
						return &storage.Final{
							Error: errors.New("插件源异常！"),
						}
					}
					zipfile := plugin_dir + "/" + utils.GenUUID() + ".zip"
					err = func() error {
						defer resp.Body.Close()
						f, err := os.OpenFile(zipfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
						if err != nil {
							return errors.New("下载异常！")
						}
						defer f.Close()
						_, err = io.Copy(f, resp.Body)
						if err != nil {
							return errors.New("文件异常！")
						}
						return nil
					}()
					if err != nil {
						return &storage.Final{
							Error: err,
						}
					}
					defer os.Remove(zipfile)
					if err := unzip(zipfile, 0755); err != nil {
						return &storage.Final{
							Error: errors.New("安装异常！"),
						}
					}
					return &storage.Final{
						Now: "",
					}
				}
				script := string(fetchScript(p.Address, key))
				if f, _, _ := initPlugin(script, p.UUID, ""); f.CreateAt != "" {
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
				break
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
							// su.DeleteValue("admin")
							su.DeleteValue("priority")
							su.DeleteValue("platform")
						}
					} else {
						su.DeleteValue("rule")
						su.DeleteValue("cron")
						su.DeleteValue("web")
						// su.DeleteValue("admin")
						su.DeleteValue("priority")
						su.DeleteValue("platform")
					}
				} else {
					// su.DeleteValue("rule")
					// su.DeleteValue("cron")
					su.DeleteValue("web")
					// su.DeleteValue("admin")
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
		f, cbs, err := initPlugin(new, key, "")
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
						console.Log("已加载 %s%s", f.Title, f.Suffix)
					} else if !f.OnStart {
						console.Log("已重载 %s%s", f.Title, f.Suffix)
					}
				} else {
					of, _, _ := initPlugin(old, key, "")
					console.Log("已卸载 %s%s", of.Title, of.Suffix)
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

func initPlugin(data string, uuid string, scriptType string) (*common.Function, []func(), error) {
	f, cbs := pluginParse(data, uuid)
	f.Suffix = ".js"
	f.Type = "goja"
	script := ""
	if f.Encrypt {
		script = DecryptPlugin(string(data))
	} else {
		script = string(data)
	}
	script = halfDeEct(script)
	script = strings.ReplaceAll(script, "new Sender", "Sender")
	script = strings.ReplaceAll(script, "new Bucket", "Bucket")
	// script = regexp.MustCompile(`import\s+\{\s*([^\}]+)\s*\}\s*from\s*['"]([^'"]+)['"]\s*;`).ReplaceAllString(script, "const {$1} = require('$2');")
	// script = regexp.MustCompile(`import\s+\s*([^\}]+)\s*\s*from\s*['"]([^'"]+)['"]\s*;`).ReplaceAllString(script, "const $1 = require('$2');")
	var err error
	prg, err2 := goja.Compile(f.Title+".js", script, false)
	if err == nil && err2 != nil {
		err = err2
	}
	// if err == nil && len(rules) == 0 && cron != "" {
	// 	err = fmt.Errorf("无效的脚本%s", title)
	// }
	// if icon == "" {
	// 	icon = "https://joeschmoe.io/api/v1/random?t=" + fmt.Sprint(time.Now().Nanosecond())
	// }
	var running func() bool

	f.Handle = func(s common.Sender, set func(vm *goja.Runtime)) interface{} {
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
			SetPluginMethod(vm, uuid, f.OnStart, running)
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
			vm.Set("action", goja.Undefined())
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
