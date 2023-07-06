package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"
	"github.com/google/uuid"
)

// type ProxyKey struct {
// }

type ProxyConfig struct {
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
		new := string(b2)
		var ncfg = ProxyConfig{}
		if strings.HasPrefix(new, "o:") {
			err := json.Unmarshal([]byte(strings.Replace(new, "o:", "", 1)), &ncfg)
			if err != nil {
				console.Log("无法解析的代理数据：", err)
				return nil
			}
			// ncfg.RuleMatcher, err = preprocessRules(ncfg.Rules)
			if err != nil {
				console.Log("无法处理的代理匹配规则：", err)
				return nil
			}
			p, err := adapter.ParseProxy(structToMap(ncfg))
			if err != nil {
				console.Log("无法解析的代理：", err)
				return nil
			}
			ncfg.Conn = p
			// fmt.Println("===", string(utils.JsonMarshal(ncfg)))
			// fmt.Println(ncfg.RuleMatcher.Match("106.52.87.206"))
			// fmt.Println(ncfg.RuleMatcher.Match("api.telegram.org"))
			Proxies.Store(string(b1), &ncfg)
		}
		return nil
	})
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
		return &storage.Final{
			Now: new,
		}
	})
}

func structToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	d, _ := json.Marshal(s)
	json.Unmarshal(d, &result)
	return result
}

func IsDifferent(a, b interface{}, ignoreFields []string) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Struct || bVal.Kind() != reflect.Struct {
		// 如果不是 struct 类型，则返回 true（不同）
		return true
	}

	// 获取 struct 的字段数
	numFields := aVal.NumField()

	// 遍历 struct 的所有字段，检查是否有不同
	for i := 0; i < numFields; i++ {
		aField := aVal.Field(i)
		bField := bVal.Field(i)

		// 获取字段名
		fieldName := aVal.Type().Field(i).Name

		// 如果字段名在 ignoreFields 中，则跳过比较
		if Contains(ignoreFields, fieldName) {
			continue
		}

		// 如果字段类型不同，则返回 true（不同）
		if aField.Type() != bField.Type() {
			return true
		}

		// 如果字段值不同，则返回 true（不同）
		if !reflect.DeepEqual(aField.Interface(), bField.Interface()) {
			return true
		}
	}

	// 如果所有字段都相同，则返回 false（相同）
	return false
}
