package core

import (
	"bytes"
	"fmt"
	"net/http"
	u "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/goccy/go-json"
)

func fetch(vm *goja.Runtime, wts ...interface{}) (interface{}, error) {
	var method = "get"
	var url = ""
	var req *http.Request

	var headers map[string]interface{}
	var transport *http.Transport
	var formData map[string]interface{}
	// var isJson bool
	var isJsonBody bool
	var body string
	var allow_redirects bool = true
	var responseType = ""
	// var err error
	// var useproxy bool
	var timeout time.Duration = 0
	var login = false
	for _, wt := range wts {
		switch wt := wt.(type) {
		case string:
			url = wt
		case map[string]interface{}:
			props := wt
			for i := range props {
				switch strings.ToLower(i) {
				case "timeout":
					timeout = time.Duration(utils.Int64(props[i])) * time.Millisecond
				case "headers":
					headers = props[i].(map[string]interface{})
				case "method":
					method = strings.ToLower(props[i].(string))
				case "url":
					if f, ok := props[i].(func(i goja.FunctionCall) goja.Value); ok {
						url = f(goja.FunctionCall{}).ToString().String()
					} else {
						url = props[i].(string)
					}
				case "json":
					if props[i].(bool) {
						responseType = "json"
					}
				case "responsetype":
					responseType = props[i].(string)
				case "datatype":
					responseType = props[i].(string)
				case "allowredirects":
					allow_redirects = props[i].(bool)
				case "body":
					if v, ok := props[i].(string); !ok {
						d, _ := json.Marshal(props[i])
						body = string(d)
						isJsonBody = true
					} else {
						body = v
					}
				case "login":
					login = true
				case "formdata":
					formData = props[i].(map[string]interface{})
				case "form":
					formData = props[i].(map[string]interface{})
				case "proxy":
					var url, user, password string
					for k, v := range props[i].(map[string]interface{}) {
						if k == "url" {
							url = fmt.Sprint(v)
						}
						if k == "user" {
							user = fmt.Sprint(v)
						}
						if k == "password" {
							password = fmt.Sprint(v)
						}
					}
					if url != "" {
						var err error
						transport, err = GetTransport(url, user, password)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
		method = strings.ToUpper(method)
		var err error
		if len(formData) > 0 {
			data := u.Values{}
			for key, value := range formData {
				switch v := value.(type) {
				case string:
					data.Set(key, v)
				case bool:
					data.Set(key, strconv.FormatBool(v))
				case int:
					data.Set(key, strconv.Itoa(v))
				// 可以根据需要添加其他类型的处理
				default:
					data.Set(key, fmt.Sprintf("%v", v))
				}
			}

			req, err = http.NewRequest(method, url, strings.NewReader(data.Encode()))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
			if err != nil {
				return nil, err
			}
			if isJsonBody {
				req.Header.Set("Content-Type", "application/json")
			}
		}
		if login {
			req.Header.Set("Cookie", "uuid=40e67d5e-f6f3-11ed-8bc2-dca9049272e5; token="+getTempAuth())
		}
		for i := range headers {
			req.Header.Set(i, fmt.Sprint(headers[i]))
		}
	}

	var rspObj goja.Proxy
	var rsp *http.Response
	var err error

	var client = &http.Client{
		Timeout: timeout,
	}
	if transport != nil {
		client.Transport = transport
	}

	if !allow_redirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	rsp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	obj, err := MakeResponseObject(vm, rsp, responseType)
	rspObj = vm.NewProxy(obj, &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			obj := target.Get(property)
			if obj != nil {
				return obj
			}
			switch property {
			case "statusText", "statusMessage":
				return vm.ToValue("")
			case "statusCode":
				return target.Get("status")
			// case "body":
			// 	return vm.ToValue(target.Get("getBody").Export().(func() interface{})())
			case "then":
				return goja.Undefined()
			}
			console.Error("response has no property " + property)
			return goja.Undefined()
		},
	})
	return rspObj, err
}
