package core

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/robfig/cron/v3"
)

type ScriptUtils struct {
	matched bool
	ress    [][]string
	script  string
}

func (su *ScriptUtils) match() {
	su.ress = regexp.MustCompile(
		`(\x20?\*[ ]?@([^\s]+)\s+([^\n]+?)\n)`,
	).FindAllStringSubmatch(su.script, -1)
	su.matched = true
}

func (su *ScriptUtils) GetValue(key string) string {
	if !su.matched {
		su.match()
	}
	value := ""
	for _, res := range su.ress {
		if res[2] == key {
			value = res[3]
		}
	}
	return value
}

// func Error(err error) interface{} {
// 	if err != nil {
// 		return map[string]string{"name": "Error", "message": err.Error()}
// 	}
// 	return nil
// }

func (su *ScriptUtils) SetValue(key, value string) {
	if !su.matched {
		su.match()
	}
	exists := []string{}
	first := ""
	for _, res := range su.ress {
		if first == "" {
			first = res[1]
		}
		if res[2] == key {
			exists = append(exists, res[1])
		}
	}
	if len(exists) != 0 {
		for i := range exists {
			if i == len(exists)-1 {
				su.script = strings.Replace(su.script, exists[i], fmt.Sprintf(" * @%s %s\n", key, value), 1)
			} else {
				su.script = strings.Replace(su.script, exists[i], "", 1)
			}
		}
	} else {
		if first != "" {
			su.script = strings.Replace(su.script, first, first+fmt.Sprintf(" * @%s %s\n", key, value), 1)
		} else {
			su.script = fmt.Sprintf("/**\n%s */\n", fmt.Sprintf(" * @%s %s\n", key, value)) + su.script
		}
	}
	su.match()
}

func (su *ScriptUtils) DeleteValue(key string) {
	if !su.matched {
		su.match()
	}
	exists := []string{}
	for _, res := range su.ress {
		if res[2] == key {
			exists = append(exists, res[1])
		}
	}
	if len(exists) != 0 {
		for i := range exists {
			su.script = strings.Replace(su.script, exists[i], "", 1)
		}
	}
	su.match()
}

var store sync.Map
var crons sync.Map

func CancelPluginCrons(uuid string) {
	v, ok := crons.Load(uuid)
	if ok {
		for _, id := range *v.(*[]cron.EntryID) {
			C.Remove(id)
		}
	}
}

func SetPluginMethod(vm *goja.Runtime, uuid string, on_start bool, running func() bool) {
	vm.Set("Bucket", func(name string) interface{} {
		return vm.NewProxy(MakeBucketObject(vm, uuid, on_start, MakeBucket(name)), &goja.ProxyTrapConfig{
			Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
				obj := target.Get(property)
				if obj != nil {
					return obj
				}
				result := target.Get("get").Export().(func(...interface{}) interface{})(property)
				return vm.ToValue(result)
			},
			Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
				result := target.Get("set").Export().(func(interface{}, interface{}) error)(
					property, value.Export(),
				)
				return result == nil
			},
		})
	})

	sillyGirlJsIplm := func(call goja.ConstructorCall) *goja.Object {
		userId := fmt.Sprintf("%d", rand.Int63())
		call.This.Set("isSlaveMode", func() bool {
			return utils.SlaveMode
		})
		call.This.Set("uuid", func() string {
			return uuid
		})
		call.This.Set("store", func(str string, v interface{}) {
			store.Store(str, v)
		})
		call.This.Set("load", func(str string) interface{} {
			v, _ := store.Load(str)
			return v
		})
		call.This.Set("delete", func(str string) {
			store.Delete(str)
		})
		call.This.Set("session", func(info interface{}) func(...int) interface{} {
			msg := ""
			imTpye := "carry"
			var chatId string
			switch info := info.(type) {
			case string:
				msg = info
			default:
				props := info.(map[string]interface{})
				for i := range props {
					switch strings.ToLower(i) {
					case "imtype":
						imTpye = props[i].(string)
					case "platform":
						imTpye = props[i].(string)
					case "msg", "message":
						msg = props[i].(string)
					case "chatid", "chat_id":
						chatId = fmt.Sprint(props[i].(string))
					case "userid", "user_id":
						userId = props[i].(string)
					}
				}
			}
			if msg == "" {
				return nil
			}
			c := &Faker{
				Type:    imTpye,
				Message: msg,
				Carry:   make(chan string),
				UserID:  userId,
				ChatID:  chatId,
				Admin:   true,
			}
			Messages <- c
			var f = func(i ...int) interface{} {
				timeOut := 1000 * 100
				if len(i) > 0 {
					timeOut = i[0]
				}
				select {
				case v, ok := <-c.Listen():
					return map[string]interface{}{
						"hasNext": ok,
						"message": v,
					}
				case <-time.After(time.Millisecond * time.Duration(timeOut)):
					return map[string]interface{}{
						"hasNext": false,
						"message": "已超时",
					}
				}
			}
			return f
		})
		return nil
	}
	registry := require.NewRegistry(require.WithLoader(mapFileSystemSourceLoader(uuid)))
	if on_start {
		var ids = &[]cron.EntryID{}
		crons.Store(uuid, ids)
		vm.Set("Cron", func() *goja.Object {
			o := vm.NewObject()
			o.Set("add", func(cron string, f func()) interface{} {
				if f == nil {
					return map[string]interface{}{
						"id":    0,
						"error": "未传入handler",
					}
				}
				cron = strings.TrimSpace(cron)
				if len(regexp.MustCompile(`\S+`).FindAllString(cron, -1)) == 5 {
					cron = "0 " + cron
				}
				id, err := C.AddFunc(cron, func() {
					// mutex := GetMutex(uuid)
					// mutex.Lock()
					// defer mutex.Unlock()
					defer func() {
						err := recover()
						if err != nil {
							console.Error("C.AddFunc err: %v", err)
						}
					}()
					f()
				})
				if err == nil {
					*ids = append(*ids, id)
				}
				return map[string]interface{}{
					"id":    id,
					"error": err,
				}
			})
			o.Set("remove", func(id cron.EntryID) {
				C.Remove(id)
			})
			return o
		})

		// registry.RegisterNativeModule("express", func(runtime *goja.Runtime, module *goja.Object) {
		// 	o := module.Get("exports").(*goja.Object)
		// 	methods := []string{"get", "post", "delete", "put", "fetch"}
		// 	for i := range methods {
		// 		method := methods[i]
		// 		o.Set(method, func(path string, handles ...func(*Request, *Response)) {
		// 			webs = append(webs, Web{
		// 				uuid:    uuid,
		// 				method:  strings.ToUpper(method),
		// 				path:    path,
		// 				handles: handles,
		// 			})
		// 		})
		// 	}
		// })
		vm.Set("gofor", func(running func() bool, handle func()) {
			go func() {
				for {
					if !func() bool {
						defer func() {
							err := recover()
							if err != nil {
								console.Error("gofor running error:", err)
							}
						}()
						return running()
					}() {
						return
					}
					func() {
						defer func() {
							err := recover()
							if err != nil {
								console.Error("gofor error:", err)
							}
						}()
						handle()
					}()
				}
			}()
		})
		// func (f *Factory) Request(running func() bool, wt interface{}, handles ...func(error, map[string]interface{}, interface{}) interface{}) {
		// 	go func() {
		// 		for {
		// 			if !running() {
		// 				return
		// 			}
		// 			func() {
		// 				defer func() {
		// 					if debug_mode {
		// 						return
		// 					}
		// 					err := recover()
		// 					if err != nil {
		// 						console.Error("Sender(\""+f.platform+"\").request error:", err)
		// 					}
		// 				}()
		// 				request(wt, handles...)
		// 			}()

		// 		}
		// 	}()
		// }
	}
	vm.Set("Express", func() *goja.Object {
		o := vm.NewObject()
		methods := []string{"get", "post", "delete", "put", "fetch"}
		for i := range methods {
			method := methods[i]
			o.Set(method, func(path string, handles ...func(*Request, *Response)) {
				webs = append(webs, Web{
					uuid:    uuid,
					method:  strings.ToUpper(method),
					path:    path,
					handles: handles,
				})
			})
		}
		o.Set("static", func(path string) {
			addStatic(uuid, path)
		})
		return o
	})
	registry.Enable(vm)
	vm.SetFieldNameMapper(myFieldNameMapper{})
	// vm.Set("sleep", sleep)
	vm.Set("md5", utils.Md5)
	vm.Set("image", utils.ToImageQrcode)
	vm.Set("video", utils.ToVideoQrcode)
	vm.Set("console", console)
	vm.Set("sillyGirl", sillyGirlJsIplm)
	vm.Set("call", func(str string) interface{} {
		return RegistFuncs[str]
	})
	vm.Set("fmt", &Fmt{})
	vm.Set("strings", &Strings{})
	// vm.Set("Bucket", BucketJsImpl)
	vm.Set("time", &TimeJsImpl{
		Second: time.Second,
		Minute: time.Minute,
		Hour:   time.Hour,
		Day:    time.Hour * 24,
	})
	vm.Set("Regexp", Regexp)
	vm.Set("Form", func(...interface{}) {

	})
	vm.Set("url2Base64", Url2Base64)
	// 字符串转Base64编码
	vm.Set("stringToBase64", stringToBase64)
	vm.Set("base64ToString", base64ToString)
	vm.Set("Buffer", func(call goja.ConstructorCall) *goja.Object {
		return Buffer(vm, call)
	})
	vm.Set("fetch", func(wts ...interface{}) interface{} {
		promise, resolve, reject := vm.NewPromise()
		func() {
			func() {
				v := recover()
				if v != nil {
					if err, ok := v.(error); ok {
						reject(Error(vm, err))
					} else {
						reject(Error(vm, fmt.Sprint(v)))
					}

				}
			}()
			fetch(vm, resolve, reject, wts...)
		}()
		return promise
	})
	vm.Set("request", func(wts ...interface{}) interface{} {
		return fetch(vm, nil, nil, wts...)
	})
	// for _, method := range []string{"get", "post", "delete", "put", "fetch"} {
	// 	vm.Set(method, )
	// }
	vm.Set("HttpListen", func(method, api string) *goja.Promise {
		promise, resolve, reject := vm.NewPromise()
		go func() {
			func() {
				v := recover()
				if v != nil {
					if err, ok := v.(error); ok {
						reject(Error(vm, err))
					} else {
						reject(Error(vm, fmt.Sprint(v)))
					}

				}
			}()
			AddHttpListen(api, strings.ToUpper(method), vm, uuid, resolve, reject)
		}()
		return promise
	})
	vm.Set("getReplyMessage", func(plt string, bots_id []string) *goja.Promise {
		return GetReplyMessage(vm, plt, bots_id)
	})
	vm.Set("Script", func(str string) interface{} {
		if str == "" {
			str = uuid
		}
		return Script(str)
	})
	vm.Set("Temp", func(pre string, sec int) interface{} {
		return map[string]interface{}{
			"set": func(key string, value interface{}, num int) {
				if sec != 0 && num == 0 {
					num = sec
				}
				temp.Set(pre+"_"+key, value, num)
			},
			"get": func(key string, def interface{}) interface{} {
				v := temp.Get(pre + "_" + key)
				if v == nil {
					v = def
				}
				return v
			},
			"delete": func(key string) {
				temp.Delete(pre + "_" + key)
			},
		}
	})

	osjs := getJsOs(vm, running)
	vm.Set("os", osjs)
	vm.Set("fs", osjs)
}

func EncryptPlugin(script string) string {
	res := strings.SplitN(script, "*/\n", 2)
	if len(res) != 2 {
		// logs.Info(len(res))
		return script
	}
	str, err := EncryptByAes([]byte(res[1]))
	if err != nil {
		// logs.Info(err)
		return script
	}
	su := ScriptUtils{script: res[0]}
	if su.GetValue("encrypt_data") != "" {
		return script
	}
	su.SetValue("encrypt_data", str)
	return su.script + "*/\n"
}

func DecryptPlugin(script string) string {
	su := ScriptUtils{script: script}
	encrypt_data := su.GetValue("encrypt_data")
	if encrypt_data == "" {
		return script
	}
	su.DeleteValue("encrypt_data")
	str, err := DecryptByAes(encrypt_data)
	if err != nil {
		return script
	}
	return fmt.Sprintf("%s%s", su.script, str)
}
