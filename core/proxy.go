package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/google/uuid"
)

// type ProxyKey struct {
// }

type ProxyConfig struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Server string `json:"server"`
	Port   int    `json:"port"`
	Name   string `json:"name"`
	// ProxyKey
	// RuleMatcher    *RuleMatcher `json:"-"`
	Conn    C.Proxy  `json:"-"`
	Rules   []string `json:"rules"`
	Plugins []string `json:"plugins"`

	// Cipher         string `json:"cipher,omitempty"`
	// Username       string `json:"username,omitempty"`
	// Password       string `json:"password,omitempty"`
	// SkipCertVerify bool   `json:"skip-cert-verify,omitempty"`
	// TLS            bool   `json:"tls,omitempty"`
	// CreatedAt int    `json:"created_at,omitempty"`
	// Remark    string `json:"remark,omitempty"`
	Enable bool `json:"enable,omitempty"`
	// UDP            bool              `json:"udp,omitempty"`
	// Plugin string `json:"plugin,omitempty"`
	// PluginOpts     map[string]string `json:"plugin-opts,omitempty"`
	// Obfs           string            `json:"obfs,omitempty"`
	// ObfsOpts       map[string]string `json:"obfs-opts,omitempty"`
	// Enable         bool     `json:"enable,omitempty"`
	// ExcludeIPs     []string `json:"exclude-ip,omitempty"`
	// ProxyGroups    []string `json:"proxy-groups,omitempty"`
	Google []int `json:"google"`
}

var Proxies sync.Map
var proxies = MakeBucket("proxies")

func GetProxyTransport(rawURL string, uuid string, params map[string]interface{}) (C.Conn, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return nil, err
		}
	}
	addr := C.Metadata{
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: port,
	}
	var p *ProxyConfig
	if params != nil {
		params["port"] = utils.Int(params["port"])
		conn, err := adapter.ParseProxy(params)
		if err != nil {
			err = fmt.Errorf("代理配置错误：%v", err)
			return nil, err
		}
		i, err := conn.DialContext(context.Background(), &addr)
		if err != nil {
			err = fmt.Errorf("代理连接错误：%v", err)
		}
		return i, err
	}
	var plugins = []*ProxyConfig{}
	Proxies.Range(func(key, value any) bool {
		cfg := value.(*ProxyConfig)
		if !cfg.Enable {
			return true
		}
		if Contains(cfg.Rules, addr.Host) {
			p = cfg
			return false
		}
		if Contains(cfg.Plugins, uuid) {
			plugins = append(plugins, cfg)
		}
		return true
	})
	if p != nil {
		i, err := p.Conn.DialContext(context.Background(), &addr)
		if err != nil {
			err = fmt.Errorf("%s(%s)代理错误：%v", p.Type, p.Name, err)
		}
		return i, err
	}
	if len(plugins) != 0 {
		p = plugins[0]
		i, err := p.Conn.DialContext(context.Background(), &addr)
		if err != nil {
			err = fmt.Errorf("%s(%s)代理错误：%v", p.Type, p.Name, err)
		}
		return i, err
	}
	return nil, nil
}

func init() {
	proxies.Foreach(func(b1, b2 []byte) error {
		key := string(b1)
		new := string(b2)
		var ncfg = ProxyConfig{}
		var params = map[string]interface{}{}
		if strings.HasPrefix(new, "o:") {
			var data = []byte(strings.Replace(new, "o:", "", 1))
			err := json.Unmarshal(data, &ncfg)
			json.Unmarshal(data, &params)
			if err != nil {
				console.Log("无法解析的代理数据：", err)
				return nil
			}
			p, err := adapter.ParseProxy(params)
			if err != nil {
				console.Log("无法解析的代理：", err)
				return nil
			}
			ncfg.Conn = p
			// fmt.Println("===", string(utils.JsonMarshal(ncfg)))
			// fmt.Println(ncfg.RuleMatcher.Match("106.52.87.206"))
			// fmt.Println(ncfg.RuleMatcher.Match("api.telegram.org"))
			Proxies.Store(key, &ncfg)
		}
		return nil
	})
	go checkProxy()
	storage.Watch(proxies, "*", func(old, new, key string) *storage.Final {
		_, err := uuid.Parse(key)
		if err != nil {
			return &storage.Final{
				Error: errors.New("非法的UUID"),
			}
		}
		// var ocfg = ProxyConfig{}
		// if strings.HasPrefix(new, "o:") {
		// 	json.Unmarshal([]byte(strings.Replace(old, "o:", "", 1)), &ocfg)
		// }
		// var okey = ProxyKey{
		// 	Type:   ocfg.Type,
		// 	Server: ocfg.Server,
		// 	Port:   ocfg.Port,
		// }
		// var nkey = ProxyKey{}
		var ncfg = ProxyConfig{}
		var params = map[string]interface{}{}
		if new == "" { //删除逻辑
			Proxies.Delete(key)
			return nil
		}

		if strings.HasPrefix(new, "o:") {
			var data = []byte(strings.Replace(new, "o:", "", 1))
			err := json.Unmarshal(data, &ncfg)
			json.Unmarshal(data, &params)
			if err != nil {
				return &storage.Final{
					Error: err,
				}
			}
			// nkey = ProxyKey{
			// 	Type:   ncfg.Type,
			// 	Server: ncfg.Server,
			// 	Port:   ncfg.Port,
			// }
		}
		// if ncfg.CreatedAt == 0 {
		// 	ncfg.CreatedAt = int(time.Now().Unix())
		// 	new = "o:" + string(utils.JsonMarshal(ncfg))
		// }
		// ov, ok := Proxies.Load(nkey)
		// if ok && (!IsDifferent(ocfg, ncfg, []string{"Name", "UUID", "Rules", "Plugins"}) || checkProxy(ov.(C.Proxy))) { //代理依旧有效
		// 	return nil
		// }
		// ncfg.RuleMatcher, err = preprocessRules(ncfg.Rules)
		if err != nil {
			return &storage.Final{
				Error: errors.New("不支持的代理匹配规则：" + err.Error()),
			}
		}
		p, err := adapter.ParseProxy(params)
		if err != nil {
			return &storage.Final{
				Error: err,
			}
		}
		ncfg.Conn = p
		Proxies.Store(key, &ncfg)
		go checkProxy()
		return &storage.Final{
			Now: new,
		}
	})
}

func checkProxy() {
	proxies.Foreach(func(b1, b2 []byte) error {
		go func(b1, b2 []byte) {
			key := string(b1)
			new := string(b2)
			var ncfg = ProxyConfig{}
			var params = map[string]interface{}{}
			if strings.HasPrefix(new, "o:") {
				var data = []byte(strings.Replace(new, "o:", "", 1))
				err := json.Unmarshal(data, &ncfg)
				json.Unmarshal(data, &params)
				if err != nil {
					return
				}
				var check = map[string]int64{}
				for _, site := range []string{"https://google.com", "https://github.com"} {
					func(site string) {
						instance, err := GetProxyTransport(site, "", params)
						if err != nil {
							return
						}
						defer instance.Close()
						var client = &http.Client{
							Transport: &http.Transport{
								Dial: func(string, string) (net.Conn, error) {
									return instance, nil
								},
								MaxIdleConns:          100,
								IdleConnTimeout:       90 * time.Second,
								TLSHandshakeTimeout:   10 * time.Second,
								ExpectContinueTimeout: 1 * time.Second,
							},
						}
						req, err := http.NewRequest("GET", site, strings.NewReader(""))
						if err != nil {
							return
						}
						startTime := time.Now().UnixMilli()
						resp, err := client.Do(req)
						if err != nil {
							return
						}
						endTime := time.Now().UnixMilli()
						resp.Body.Close()
						spend := endTime - startTime
						check[site] = spend
						// console.Log(site, "resp.StatusCode", resp.StatusCode, "time", spend)
					}(site)
				}
				params["check"] = check
				proxies.Set2(key, "o:"+string(utils.JsonMarshal(params)))
			}
		}(b1, b2)
		return nil
	})
}
