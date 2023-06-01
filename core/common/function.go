package common

type Function struct {
	Rules       []string                   `json:"-"`
	Params      [][]string                 `json:"-"`
	ImType      *Filter                    `json:"-"`
	UserId      *Filter                    `json:"-"`
	GroupId     *Filter                    `json:"-"`
	FindAll     bool                       `json:"-"`
	Admin       bool                       `json:"-"`
	Handle      func(s Sender) interface{} `json:"-"`
	Cron        string                     `json:"cron"`
	Priority    int                        `json:"-"`
	Disable     bool                       `json:"-"`
	Hidden      bool                       `json:"-"`
	CronId      int                        `json:"-"`
	Origin      string                     `json:"-"`
	UUID        string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Public      bool                       `json:"public"`
	Icon        string                     `json:"icon"`
	Version     string                     `json:"version"`
	Author      string                     `json:"author"`
	Status      int                        `json:"status"` //0未安装 1可更新 2已安装
	Address     string                     `json:"-"`
	CreateAt    string                     `json:"create_at"`
	Module      bool                       `json:"module"`
	// Web         bool                       `json:"web"`
	Encrypt bool `json:"encrypt"`
	OnStart bool `json:"on_start"`
	PluginPublisher
	Running bool `json:"-"`
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
