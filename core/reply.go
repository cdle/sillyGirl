package core

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/buger/jsonparser"
)

type Reply struct {
	Rules   []string
	Type    string //text url
	Content string
	Request struct {
		Url          string
		Method       string
		Body         string // form-data raw
		Headers      []string
		ResponseType string `yaml:"response_type"` //text json image
		Get          string
		Regex        string
		Template     string
	}
	Replace [][]string
}

func InitReplies() {
	for _, v := range Config.Replies {
		reply := v
		var handler func(s Sender) interface{}
		if reply.Type != "url" {
			handler = func(s Sender) interface{} {
				return reply.Content
			}
		}
		handler = func(s Sender) interface{} {
			url := reply.Request.Url
			body := reply.Request.Body
			for k, v := range s.GetMatch() {
				url = strings.Replace(url, fmt.Sprintf(`{{%d}}`, k+1), v, -1)
				body = strings.Replace(body, fmt.Sprintf(`{{%d}}`, k+1), v, -1)
			}
			var req httplib.BeegoHTTPRequest
			if strings.ToLower(reply.Request.Method) == "post" {
				req = *httplib.Post(url)
			} else {
				req = *httplib.Get(url)
			}
			for _, header := range reply.Request.Headers {
				ss := strings.Split(header, ":")
				if len(ss) > 0 {
					req.Header(ss[0], strings.Join(ss[1:], ":"))
				}
			}
			if reply.Request.Body != "" {
				req.Body(body)
			}
			rsp, err := req.Response()
			if err != nil {
				if reply.Content != "" {
					s.Reply(reply.Content)
				} else {
					s.Reply(err)
				}
				return nil
			}
			switch reply.Request.ResponseType {
			case "image":
				if reply.Request.Get != "" {
					d, _ := ioutil.ReadAll(rsp.Body)
					f, err := jsonparser.GetString(d, strings.Split(reply.Request.Get, ".")...)
					if err != nil {
						s.Reply(err)
						return nil
					}
					s.Reply(httplib.Get(f).Response())
					return nil
				}
				if reply.Request.Regex != "" {
					d, _ := ioutil.ReadAll(rsp.Body)
					res := regexp.MustCompile(reply.Request.Regex).FindStringSubmatch(string(d))
					if len(res) != 0 {
						s.Reply(httplib.Get(res[1]).Response())
					}
					return nil
				}
				s.Reply(rsp)
			case "json":
				d, _ := ioutil.ReadAll(rsp.Body)
				f, err := jsonparser.GetString(d, strings.Split(reply.Request.Get, ".")...)
				if err != nil {
					s.Reply(err)
					return true
				}
				s.Reply(f)
			case "template":
				data, _ := ioutil.ReadAll(rsp.Body)
				content := reply.Request.Template
				for _, re := range regexp.MustCompile(`gjson[(][^()]+[)]`).FindAllStringSubmatch(content, -1) {
					v := re[0]
					get := strings.Replace(strings.TrimRight(v, ")"), "gjson(", "", -1)
					f, _ := jsonparser.GetString(data, strings.Split(get, ".")...)
					content = strings.Replace(content, v, f, -1)
				}
				for i := range reply.Replace {
					if len(reply.Replace[i]) >= 2 {
						content = strings.Replace(content, reply.Replace[i][0], reply.Replace[i][1], -1)
					} else {
						content = strings.Replace(content, reply.Replace[i][0], "", -1)
					}
				}
				for _, re := range regexp.MustCompile(`image[(][^()]+[)]`).FindAllStringSubmatch(content, -1) {
					v := re[0]
					get := strings.Replace(strings.TrimRight(v, ")"), "image(", "", -1)
					f, _ := jsonparser.GetString(data, strings.Split(get, ".")...)
					s.Reply(httplib.Get(f).Response())
					content = strings.Replace(content, v, "", -1)
				}
				s.Reply(content)
			default:
				d, _ := ioutil.ReadAll(rsp.Body)
				s.Reply(d)
			}
			return nil
		}
		functions = append(functions, Function{
			Rules:  reply.Rules,
			Handle: handler,
		})

	}
}
