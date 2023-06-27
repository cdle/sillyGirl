package core

import (
	"io/ioutil"
	"time"

	"github.com/cdle/sillyGirl/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

var web_sessions = MakeBucket("web_sessions")

type Web struct {
	uuid    string
	method  string
	path    string
	handles []func(*Request, *Response)
}

var webs []Web

func CancelPluginWebs(uuid string) {
	for i := range webs {
		if webs[i].uuid == uuid {
			webs[i].handles = nil
		}
	}
}

type Request struct {
	c          *gin.Context
	ress       [][]string
	parsedForm bool
	handled    bool
	uuid       string
	bodyData   []byte
	_event     string
}

func (r *Request) Body() string {
	if len(r.bodyData) == 0 {
		r.bodyData, _ = ioutil.ReadAll(r.c.Request.Body)
	}
	return string(r.bodyData)
}

func (r *Request) Json() interface{} {
	var i interface{}
	if len(r.bodyData) == 0 {
		r.bodyData, _ = ioutil.ReadAll(r.c.Request.Body)
	}
	if json.Unmarshal(r.bodyData, &i) != nil {
		return nil
	}
	return i
}

func (r *Request) Ip() string {
	return r.c.ClientIP()
}

func (r *Request) Event() string {
	return r._event
}

func (r *Request) OriginalUrl() string {
	return r.c.Request.URL.String()
}

func (r *Request) Query(param string) string {
	return r.c.Query(param)
}

func (r *Request) Param(i int) string {
	return r.ress[i-1][1]
}

func (r *Request) Querys() map[string][]string {
	return r.c.Request.URL.Query()
}

func (r *Request) PostForm(s string) string {
	if !r.parsedForm {
		r.c.Request.ParseForm()
	}
	return r.c.PostForm(s)
}

func (r *Request) PostForms() map[string][]string {
	if !r.parsedForm {
		r.c.Request.ParseForm()
	}
	return r.c.Request.PostForm
}

func (r *Request) Path() string {
	return r.c.Request.URL.Path
}

func (r *Request) Header(s string) string {
	return r.c.GetHeader(s)
}

func (r *Request) Get(s string) string {
	return r.c.GetHeader(s)
}

func (r *Request) Is(s string) bool { //判断是否传入的MIME类型
	return true
}

func (r *Request) Headers() map[string][]string {
	return r.c.Request.Header
}

func (r *Request) Method() string {
	return r.c.Request.Method
}

func (r *Request) Logined() bool {
	auth, _ := CheckAuth(r.Cookie("token"))
	return auth != nil
}

func (r *Request) Cookie(s string) string {
	var cookie, _ = r.c.Cookie(s)
	return cookie
}

func (r *Request) Cookies() map[string]string {
	var cookies = map[string]string{}
	for _, v := range r.c.Request.Cookies() {
		cookies[v.Name] = v.Value
	}
	return cookies
}

func (r *Request) Continue() {
	r.handled = false
}

func (r *Request) SetSession(k, v string) string {
	j := map[string]interface{}{}
	json.Unmarshal(web_sessions.GetBytes(r.uuid), &j)
	j[k] = v
	j["time"] = time.Now().Unix()
	_, err := web_sessions.Set(r.uuid, utils.JsonMarshal(j))
	if err != nil {
		return err.Error()
	}
	return ""
}

func (r *Request) GetSession(k string) string {
	j := map[string]interface{}{}
	json.Unmarshal(web_sessions.GetBytes(r.uuid), &j)
	v, ok := j[k].(string)
	if !ok {
		return ""
	}
	return v
}

func (r *Request) GetSessionID() string {
	return r.uuid
}

func (r *Request) DestroySession() string {
	return ""
}
