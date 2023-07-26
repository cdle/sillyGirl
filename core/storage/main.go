package storage

type Listen struct {
	UUID   string
	Name   string
	Key    string
	Handle func(old string, new string, key string) *Final
}

const EMPTY = "__EMPTY__"

var Listens []Listen

var DisableHandle = func(uuid string) {
	for i := range Listens {
		if Listens[i].UUID == uuid {
			Listens[i].Handle = func(old, new, key string) *Final {
				return nil
			}
		}
	}
}

type Bucket interface {
	Set(interface{}, interface{}) (string, bool, error)
	Set2(interface{}, interface{}) (string, bool, error)
	Copy(string) Bucket
	IsEmpty() (bool, error)
	Size() (int64, error)
	Delete() error
	Type() string
	Buckets() []string
	GetString(...interface{}) string
	GetBytes(string) []byte
	GetInt(string, ...int) int
	GetBool(string, ...bool) bool
	Foreach(func([]byte, []byte) error)
	Create(interface{}) error
	First(interface{}) error
	String() string
	GetName() string
	Keys() ([]string, error)
}

type Final struct {
	Now     string
	Message string
	Error   error
	EndFunc func()
}

func Watch(bucket Bucket, key interface{}, handle func(old string, new string, key string) *Final, uuid ...string) {
	k := "*"
	if key != nil {
		k = key.(string)
	}
	listen := Listen{
		Name:   bucket.GetName(),
		Key:    k,
		Handle: handle,
	}
	if len(uuid) != 0 {
		listen.UUID = uuid[0]
	}
	Listens = append(Listens, listen)
}
