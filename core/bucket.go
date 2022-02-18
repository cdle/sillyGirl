package core

import "github.com/beego/beego/v2/adapter/logs"

var sillyGirl Bucket
var Zero Bucket

func MakeBucket(name string) Bucket {
	if Zero == nil {
		logs.Error("找不到存储器。")
	}
	return Zero.Copy(name)
}

type Bucket interface {
	Copy(name string) Bucket
	Set(key interface{}, value interface{}) error
	GetString(kv ...interface{}) string
	GetBytes(key string) []byte
	GetInt(key interface{}, vs ...int) int
	GetBool(key interface{}, vs ...bool) bool
	Foreach(f func(k, v []byte) error)
	Create(i interface{}) error
	First(i interface{}) error
	String() string
}
