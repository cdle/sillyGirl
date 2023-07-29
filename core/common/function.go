package common

import "github.com/dop251/goja"

type Function struct {
	Rules       []string                                         `json:"-"`
	Params      [][]string                                       `json:"-"`
	ImType      *Filter                                          `json:"-"`
	UserId      *Filter                                          `json:"-"`
	GroupId     *Filter                                          `json:"-"`
	FindAll     bool                                             `json:"-"`
	Admin       bool                                             `json:"-"`
	Handle      func(Sender, func(vm *goja.Runtime)) interface{} `json:"-"`
	Cron        map[string]string                                `json:"cron"`
	Priority    int                                              `json:"-"`
	Disable     bool                                             `json:"disable"`
	Hidden      bool                                             `json:"-"`
	CronIds     []int                                            `json:"-"`
	Origin      string                                           `json:"-"`
	UUID        string                                           `json:"id"`
	Title       string                                           `json:"title"`
	Type        string                                           `json:"type"`   //脚本类型
	Suffix      string                                           `json:"suffix"` //脚本后缀
	Description string                                           `json:"description"`
	Public      bool                                             `json:"public"`
	Icon        string                                           `json:"icon"`
	Version     string                                           `json:"version"`
	Author      string                                           `json:"author"`
	Status      int                                              `json:"status"` //0未安装 1可更新 2已安装
	Address     string                                           `json:"-"`
	CreateAt    string                                           `json:"create_at"`
	Module      bool                                             `json:"module"`
	// Web         bool                       `json:"web"`
	Encrypt bool `json:"encrypt"`
	OnStart bool `json:"on_start"`
	PluginPublisher
	Running   bool        `json:"running"`
	Https     []*Http     `json:"-"`
	Reply     *Reply      `json:"-"`
	Downloads int         `json:"downloads"`
	HasForm   bool        `json:"has_form"`
	Carry     bool        `json:"carry"`
	Messages  interface{} `json:"messages"`
	Classes   []string    `json:"classes"`
	Debug     bool        `json:"debug"`
	Path      string      `json:"-"`
	Reload    func()      `json:"-"`
}
type Filter struct {
	BlackMode bool
	Items     []string
}

type PluginPublisher struct {
	Address      string `json:"address"`
	Organization string `json:"organization"`
	Identified   bool   `json:"identified"`
}

type Http struct { //GET /abc
	Path   string
	Method string
}

type Reply struct { //wx 123 1234
	Platform string
	BotsID   []string
}
