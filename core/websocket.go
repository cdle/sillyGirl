package core

import (
	"net/http"
	"strings"

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
					var res = &Response{
						c: c,
					}
					req._event = "connect"
					function.Handle(&Faker{
						Type: "*",
					}, func(vm *goja.Runtime) {
						vm.Set("res", res)
						vm.Set("req", req)
					})
					// message
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
					res.conn = ws
					for {
						_, data, err := ws.ReadMessage()
						if err != nil { // disconnect
							req._event = "disconnect"
							for _, f2 := range Functions {
								if f2.UUID == function.UUID {
									function = f2
								}
							}
							function.Handle(&Faker{
								Type: "*",
							}, func(vm *goja.Runtime) {
								vm.Set("res", res)
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
						function.Handle(&Faker{
							Type: "*",
						}, func(vm *goja.Runtime) {
							vm.Set("res", res)
							vm.Set("req", req)
						})
					}
					return
				}
			}
		}
	}
}
