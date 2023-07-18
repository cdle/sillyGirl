package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/goccy/go-json"
)

func fetch(vm *goja.Runtime, uuid string, wts ...interface{}) (interface{}, error) {
	var method = "get"
	var url = ""
	var req *http.Request

	var headers map[string]interface{}
	var transport *http.Transport
	var formData map[string]interface{}
	// var isJson bool
	var isJsonBody bool
	var body []byte
	var allow_redirects bool = true
	var responseType = ""
	// var err error
	// var useproxy bool
	var timeout time.Duration = 0
	var login = false
	var instance C.Conn
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
						if props[i] == nil {
							panic(Error(vm, "无效的url请求地址nil"))
						}
						url = fmt.Sprint(props[i])
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
					switch v := props[i].(type) {
					case string:
						body = []byte(v)
					case []byte:
						body = v
					case *Buffer:
						body = v.value
					default:
						d, _ := json.Marshal(props[i])
						body = d
						isJsonBody = true
					}
				case "login":
					login = true
				case "formdata":
					formData = props[i].(map[string]interface{})
				case "form":
					formData = props[i].(map[string]interface{})
				case "proxy":
					var err error
					var params = props[i].(map[string]interface{})
					if _, ok := params["name"]; !ok {
						params["name"] = "临时"
					}
					instance, err = GetProxyTransport(url, uuid, params)
					if err != nil {
						panic(Error(vm, err))
					}
					if instance != nil {
						defer instance.Close()
					}
				}
			}
		}
		var err error
		if instance == nil {
			instance, err = GetProxyTransport(url, uuid, nil)
			if err != nil {
				panic(Error(vm, err))
			}
			if instance != nil {
				defer instance.Close()
			}
		}
		if instance != nil {
			transport = &http.Transport{
				Dial: func(string, string) (net.Conn, error) {
					return instance, nil
				},
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
			}
		}
		method = strings.ToUpper(method)
		if len(formData) > 0 {
			// 创建一个新的 buffer
			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			func() {
				defer writer.Close()
				// defer func() {
				// 	ch <- true
				// }()
				for key := range formData {
					value := fmt.Sprint(formData[key])
					// 添加文本上传字段
					if !strings.HasPrefix(value, "url(") || !strings.HasSuffix(value, ")") {
						fieldWriter, err := writer.CreateFormField(key)
						if err != nil {
							console.Error("cff", err)
							return
						}
						if _, err := fieldWriter.Write([]byte(value)); err != nil {
							console.Error("fw", err)
							return
						}
						continue
					} else {
						url := strings.TrimPrefix(value, "url(")
						url = strings.TrimSuffix(url, ")")
						resp, err := http.Get(url)
						if err != nil {
							// console.Error("failed to download file %s: %v\n", url, err)
							panic(Error(vm, "failed to download file %s: %v\n", url, err))
							// return
						}
						defer resp.Body.Close()
						part, err := writer.CreateFormFile(key, filepath.Base(url))
						if err != nil {
							panic(Error(vm, "failed to get file info for %s: %v\n", url, err))
						}
						// 复制管道中的内容到表单字段中
						if _, err = io.Copy(part, resp.Body); err != nil {
							panic(Error(vm, err))
						}
					}
				}
			}()
			// 关闭表单写入器
			req, err = http.NewRequest(method, url, payload)
			// <-ch
			// ch = nil
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())
		} else {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(body))

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
