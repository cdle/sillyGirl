package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func IsWebSocketRequest(req *http.Request) bool {
	if req.Header.Get("Upgrade") != "websocket" {
		return false
	}
	if !websocket.IsWebSocketUpgrade(req) {
		return false
	}
	return true
}

type WsPattern struct {
	Value map[string]interface{}
	Chan  chan map[string]interface{}
}

type WsConn struct {
	conn     *websocket.Conn
	patterns sync.Map
	Key      int64
	sync.RWMutex
}

func (wc *WsConn) Close() error {
	return wc.conn.Close()
}

func (wc *WsConn) WriteMessage(messageType int, data []byte, pattern map[string]interface{}) (error, map[string]interface{}) {
	var res map[string]interface{}
	wp := &WsPattern{}
	var timeout int
	if pattern != nil {
		if v, ok := pattern["$timeout"]; ok {
			timeout = utils.Int(v)
			delete(pattern, "$timeout")
		}
		wp.Value = pattern
		key := atomic.AddInt64(&wc.Key, 1)
		wp.Chan = make(chan map[string]interface{}, 1)
		defer func() {
			close(wp.Chan)
			wc.patterns.Delete(key)
		}()
		wc.patterns.Store(key, wp)
	}
	var err error
	func() {
		wc.Lock()
		defer wc.Unlock()
		err = wc.conn.WriteMessage(messageType, data)
	}()
	if pattern != nil {
		if timeout == 0 {
			timeout = 5000
		}
		select {
		case res = <-wp.Chan:
		case <-time.After(time.Millisecond * time.Duration(timeout)):
		}
	}
	return err, res
}

// func (wc *WsConn) WriteMessage(messageType int, data []byte) error {
// 	return wc.conn.WriteMessage(messageType, data)
// }

func handleWebsocket(c *gin.Context) {
	for _, function := range Functions {
		if len(function.Https) != 0 {
			for _, h := range function.Https {
				path := h.Path
				method := h.Method
				if c.Request.URL.Path == path && strings.HasPrefix(method, "W") {
					// connect
					var req = &Request{
						c: c,
						// uuid: uuid,
					}

					var upGrader = websocket.Upgrader{
						CheckOrigin: func(r *http.Request) bool {
							return true
						},
					}
					ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						c.Writer.Write([]byte(err.Error()))
						return
					}
					wc := &WsConn{}

					req._event = "connect"
					wc.conn = ws
					go function.Handle(&Faker{
						Type: "websocket",
					}, func(vm *goja.Runtime) {
						vm.Set("res", &Response{
							c:    c,
							conn: wc,
							vm:   vm,
						})
						vm.Set("req", req)
					})
					time.Sleep(time.Millisecond * 500)
					for {
						_, data, err := ws.ReadMessage()
						wc.patterns.Range(func(key, value any) bool {
							wp := value.(*WsPattern)
							matched := false
							// fmt.Println("wp.Value", wp.Value)
							for k, v := range wp.Value {
								value, _, _, err := jsonparser.Get(data, strings.Split(k, ".")...)
								// fmt.Println("k, v", k, v, "key path", strings.Split(k, "."), "data:", string(data), err)
								if err != nil {
									// fmt.Println("err1", err)
									return true
								}
								if string(value) != fmt.Sprint(v) {
									return true
								}
								matched = true
								// v2 := v
								// fmt.Println("v2 := v", v2, string(value))
								// err = json.Unmarshal(value, &v2)
								// if err != nil {
								// 	fmt.Println("err2", err)
								// 	return true
								// }
								// if v != v2 {
								// 	fmt.Println("err3", err)
								// 	return true
								// } else {
								// 	matched = true
								// }
							}
							// fmt.Println("matched", matched)
							if matched {
								var result = map[string]interface{}{}
								err := json.Unmarshal(data, &result)
								if err == nil {
									select {
									case wp.Chan <- result:
									default:
									}
								} else {
									// fmt.Println("err3", err)
								}
							}
							return true
						})
						if err != nil { // disconnect
							req._event = "disconnect"
							for _, f2 := range Functions {
								if f2.UUID == function.UUID {
									function = f2
								}
							}
							function.Handle(&Faker{
								Type: "websocket",
							}, func(vm *goja.Runtime) {
								vm.Set("res", &Response{
									c:    c,
									conn: wc,
									vm:   vm,
								})
								vm.Set("req", req)
							})
							ws.Close()
							break
						}
						req.bodyData = data
						req._event = "message"
						for _, f2 := range Functions {
							if f2.UUID == function.UUID {
								function = f2
							}
						}
						go function.Handle(&Faker{
							Type: "websocket",
						}, func(vm *goja.Runtime) {
							vm.Set("res", &Response{
								c:    c,
								conn: wc,
								vm:   vm,
							})
							vm.Set("req", req)
						})
					}
					return
				}
			}
		}
	}
}
