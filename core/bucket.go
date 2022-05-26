package core

import "github.com/beego/beego/v2/adapter/logs"

var sillyGirl Bucket
var Zero Bucket

func MakeBucket(name string) Bucket {
	if Zero == nil {
		logs.Error("找不到存储器，开发者自行实现接口。")
	}
	return Zero.Copy(name)
}

type Bucket interface {
	Copy(string) Bucket
	Set(interface{}, interface{}) error
	Empty() (bool, error)
	Size() (int64, error)
	Delete() error
	Buckets() ([][]byte, error)
	GetString(...interface{}) string
	GetBytes(string) []byte
	GetInt(interface{}, ...int) int
	GetBool(interface{}, ...bool) bool
	Foreach(func([]byte, []byte) error)
	Create(interface{}) error
	First(interface{}) error
	String() string
}
