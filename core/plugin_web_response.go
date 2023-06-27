package core

import (
	"fmt"
	"net/http"

	"github.com/cdle/sillyGirl/utils"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

type Response struct {
	// Send       func(goja.Value)                     `json:"send"`
	// SendStatus func(int)                            `json:"sendStatus"`
	// Json       func(...interface{})                 `json:"json"`
	// Header     func(string, string)                 `json:"header"`
	// Render     func(string, map[string]interface{}) `json:"render"`
	// Redirect   func(...interface{})                 `json:"redirect"`
	// Status     func(int, ...string) goja.Value      `json:"status"`
	// GetStatus func() int `json:"getStatus"`
	// IsComplete func() bool `json:"isComplete"`
	// SetCookie  func(string, string, ...interface{}) `json:"setCookie"`
	c          *gin.Context
	content    string
	isJson     bool
	status     int
	isRedirect bool
	conn       *websocket.Conn
}

func (r *Response) Send(gv goja.Value) *Response {

	gve := gv.Export()
	switch gve := gve.(type) {
	case string:
		r.content += gve
	default:
		d, err := json.Marshal(gve)
		if err == nil {
			r.content += string(d)
			r.isJson = true
		} else {
			r.content += fmt.Sprint(gve)
		}
	}
	if r.conn != nil {
		r.conn.WriteMessage(1, []byte(r.content))
		r.content = ""
		return r
	}
	return r
}

func (r *Response) SendStatus(st int) *Response {
	r.status = st
	return r
}

func (r *Response) Json(ps ...interface{}) *Response {
	var status int64 = 200
	var data interface{}
	if len(ps) == 1 {
		data = ps[0]
	} else {
		status = ps[0].(int64)
		data = ps[1]
	}
	d, err := json.Marshal(data)
	f := string(d)
	r.status = int(status)
	if err == nil {
		r.content = f
		r.isJson = true
	} else {
		r.content += fmt.Sprint(data)
	}
	if r.conn != nil {
		r.conn.WriteMessage(1, d)
		r.content = ""
		return r
	}
	return r
}

func (r *Response) Header(str, value string) *Response {
	r.c.Header(str, value)
	return r
}

func (r *Response) Set(str, value string) *Response {
	r.c.Header(str, value)
	return r
}

func (r *Response) Render(path string, obj map[string]interface{}) *Response {
	r.c.HTML(http.StatusOK, path, obj)
	return r
}

func (r *Response) Redirect(is ...interface{}) {
	a := 302
	b := ""
	for _, i := range is {
		switch i := i.(type) {
		case string:
			b = i
		default:
			a = utils.Int(i)
		}
	}
	r.c.Redirect(a, b)
	r.isRedirect = true
}

func (r *Response) Status(i int, s ...string) *Response {
	r.status = i
	if len(s) != 0 {
		for _, v := range s {
			r.content += v
		}
		panic("stop")
	}
	return r
}

func (r *Response) SetCookie(name, value string, i ...interface{}) *Response {
	r.c.SetCookie(name, value, 8640000, "/", "", false, true)
	return r
}

func (r *Response) Stop() {
	panic("stop")
}
