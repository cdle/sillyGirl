package core

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/goccy/go-json"
)

// type Reason map[string]interface{}

// func readEvent(body io.Reader) (string, error) {
// 	var data []rune
// 	bufReader := bufio.NewReader(body)
// 	for {
// 		r, _, err := bufReader.ReadRune()
// 		if err != nil {
// 			return "", err
// 		}
// 		if r != '\n' {
// 			data = append(data, r)
// 		} else {

// 			return strings.Replace(string(data), "data: ", "", 1), nil
// 		}
// 	}
// 	panic(string(data))
// }

func MakeResponseObject(vm *goja.Runtime, resp *http.Response, responseType string) (*goja.Object, error) {
	obj := vm.NewObject()
	var data []byte
	var err error
	obj.Set("status", resp.StatusCode)
	obj.Set("headers", vm.NewProxy(MakeHeadersObject(vm, resp.Header), &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			obj := target.Get(property)
			if obj != nil {
				return obj
			}
			result := target.Get("get").Export().(func(name string) string)(property)
			return vm.ToValue(result)
		},
		Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
			target.Get("set").Export().(func(name, value string))(
				property, value.String(),
			)
			return true
		},
	}))
	if resp.Header.Get("Content-Type") == "text/event-stream" {
		var evc = make(chan interface{}, 100)
		var closed = false
		go func() {
			defer resp.Body.Close()
			// var data2 []rune
			var data []rune
			bufReader := bufio.NewReader(resp.Body)
			for {
				r, _, err := bufReader.ReadRune()
				if err != nil {
					evc <- err
					break
				}
				// data2 = append(data2, r)
				if r != '\n' {
					data = append(data, r)
				} else {
					event := strings.TrimSpace(strings.Replace(string(data), "data: ", "", 1))
					if event != "" {
						evc <- event
					}
					data = nil
				}
			}
			// fmt.Println(string(data2))
			time.Sleep(time.Minute * 5)
			if !closed {
				close(evc)
				evc = nil
				closed = true
			}
		}()
		obj.Set("on", func(event string) interface{} {
			switch event {
			case "data", "json":
				var rs = <-evc
				switch v := rs.(type) {
				case error:
					if !closed {
						close(evc)
						evc = nil
						closed = true
					}
					panic(Error(vm, v))
				case string:
					if event == "json" {
						var r interface{}
						err := json.Unmarshal([]byte([]byte(v)), &r)
						if err != nil {
							panic(Error(vm, err))
						}
						return r
					}
					return v
				}
			}
			return nil
		})
		return obj, err
	} else {
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return obj, err
		}
	}
	var body interface{}
	if Contains([]string{"blob", "arraybuffer"}, responseType) {
		body = data
	} else if Contains([]string{"text", "document"}, responseType) {
		body = string(data)
	} else if responseType == "json" {
		var v interface{}
		err := json.Unmarshal(data, &v)
		if err != nil { /////
			return obj, errors.New("请求返回数据不是json格式：" + string(data))
		} else {
			body = v
		}
	}

	if body == nil {
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/") {
			body = string(data)
		} else if strings.HasPrefix(contentType, "image/") {
			body = data
		}
	}

	if body == nil {
		isBinary := false
		for _, b := range data {
			if b < 32 || b > 126 {
				isBinary = true
				break
			}
		}
		if isBinary {
			body = data
		} else {
			body = string(data)
		}
	}

	obj.Set("body", body)
	obj.Set("getBody", func() interface{} {
		return body
	})
	obj.Set("json", func() interface{} {
		var v interface{}
		json.Unmarshal(data, &v)
		return v
	})
	obj.Set("text", func() string {
		return string(data)
	})

	return obj, nil
}
