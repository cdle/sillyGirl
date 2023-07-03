package core

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/goccy/go-json"
)

// type Reason map[string]interface{}

func MakeResponseObject(vm *goja.Runtime, resp *http.Response, responseType string) (*goja.Object, error) {
	obj := vm.NewObject()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil { ///////
		return obj, err
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
	return obj, nil
}
