package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/utils"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var master = ""

var sillyGirl *Bucket

type Bucket struct {
	name string
}

var toMaster = func(bucket, key, value string) (string, error) {
	req := httplib.Put(master + "/api/storage")
	web_token := sillyGirl.GetString("web_token")
	if web_token == "" {
		return "", errors.New("请先在本体可视化登录！")
	}
	token := strings.Split(web_token, "&")[0]
	req.SetCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	key = fmt.Sprintf("%s.%s", bucket, key)
	req.JSONBody(map[string]interface{}{
		key: value,
	})
	data, err := req.Bytes()
	if err != nil {
		return "", errors.New("啊，与本体失联了~")
	}
	var message string
	message, _ = jsonparser.GetString(data, "messages", key)
	errStr, _ := jsonparser.GetString(data, "errors", key)
	if errStr != "" {
		err = errors.New(errStr)
	}
	return message, err
}

var db *redis.Client

func (b *Bucket) String() string {
	return b.name
}

var Buckets = []Bucket{}

func (bucket *Bucket) Type() string {
	return "redis"
}

func Try(RedisAddr, RedisPassword string) error {
	db := redis.NewClient(&redis.Options{
		Addr:        RedisAddr,
		Password:    RedisPassword, // no password set
		DB:          0,             // use default DB
		DialTimeout: time.Second * 5,
	})
	err := db.HSet(context.Background(), "sillyGirl", "storage", "redis").Err()
	if err == nil {
		db.Close()
	}
	return err
}

func InitsillyGirl(RedisAddr, RedisPassword string) storage.Bucket {
	db = redis.NewClient(&redis.Options{
		Addr:        RedisAddr,
		Password:    RedisPassword, // no password set
		DB:          0,             // use default DB
		DialTimeout: time.Second * 5,
	})

	err := db.HSet(context.Background(), "sillyGirl", "storage", "redis").Err()
	if err != nil {
		logs.Error("redis错误", err)
		panic(err)
	}
	sillyGirl = &Bucket{
		name: "sillyGirl",
	}

	port := sillyGirl.GetString("port", "8080")
	if utils.SlaveMode {
		is, _ := db.ConfigGet(ctx, "slaveof").Result()
		db.ConfigSet(ctx, "notify-keyspace-events", "KEh")
		if len(is) > 1 {
			ip := strings.Split(fmt.Sprint(is[1]), " ")[0]
			master = fmt.Sprintf("http://%s:%s", ip, port)
		}
		go func() {
		again:
			subscriber := db.Subscribe(ctx, "__keyevent@0__:hset")
			for {
				msg, err := subscriber.ReceiveMessage(ctx)
				if err != nil {
					time.Sleep(time.Second)
					goto again
				}
				bk := msg.Payload

				// logs.Info(bk, msg.Pattern, msg.Channel, msg.Payload)
				if bk := strings.Split(bk, "."); len(bk) == 2 {
					bucket, key := bk[0], bk[1]
					for _, listen := range storage.Listens {
						if listen.Name == bucket && (listen.Key == key || listen.Key == "*") {
							str, err := db.HGet(ctx, bucket, key).Result()
							if err == nil {
								listen.Handle("", str, key)
							}
						}
					}
				}
			}
		}()
	}
	return sillyGirl
}

func (bucket *Bucket) Copy(name string) storage.Bucket {
	return &Bucket{name: name}
}

func (bucket *Bucket) Buckets() []string {
	var r []string
	r, _ = db.Keys(ctx, "*").Result()
	var e []string
	for _, v := range r {
		if strings.HasSuffix(v, "_Sequence") {
			continue
		}
		if regexp.MustCompile(`^\d+@\d+$`).FindString(v) != "" {
			continue
		}
		e = append(e, v)
	}
	return e
}

func (bucket *Bucket) Delete() error {
	db.Del(ctx, bucket.name)
	return nil
}

func (bucket *Bucket) GetName() string {
	return bucket.name
}

func Set(key string, value string, expiration time.Duration) error {
	db.Set(ctx, key, value, expiration)
	return nil
}

func Get(key string) string {
	v, _ := db.Get(ctx, key).Result()
	return v
}

func (bucket *Bucket) Set2(key interface{}, value interface{}) (string, error) {
	new := ""
	msg := ""
	k := fmt.Sprint(key)
	switch value := value.(type) {
	case []byte:
		new = string(value)
	case string:
		new = value
	case nil:
	default:
		new = fmt.Sprint(value)
	}
	if !utils.SlaveMode {
		if new == "" {
			return msg, db.HDel(ctx, bucket.name, k).Err()
		} else {
			return msg, db.HSet(ctx, bucket.name, k, new).Err()
		}
	} else {
		return toMaster(bucket.name, k, new)
	}

}

func (bucket *Bucket) Set(key interface{}, value interface{}) (string, error) {
	new := ""
	msg := ""
	k := fmt.Sprint(key)
	switch value := value.(type) {
	case []byte:
		new = string(value)
	case string:
		new = value
	case nil:
	default:
		new = fmt.Sprint(value)
	}
	if !utils.SlaveMode {
		var handles []func(string, string, string) *storage.Final
		var endFuncs = []func(){}
		for _, listen := range storage.Listens {
			if listen.Name == bucket.name && (listen.Key == key || listen.Key == "*") {
				handles = append(handles, listen.Handle)
			}
		}
		if len(handles) > 0 {
			old := bucket.GetString(key)
			if old == new {
				return msg, nil
			}
			for _, handle := range handles {
				fin := handle(old, new, k)
				if fin != nil {
					if fin.Message != "" {
						msg = fin.Message
					}
					if fin.Error != nil {
						return msg, fin.Error
					}
					if fin.Now != "" {
						new = fin.Now
					}
					if fin.EndFunc != nil {
						endFuncs = append(endFuncs, fin.EndFunc)
					}
				}
			}
		}
		var err error
		if new == "" {
			err = db.HDel(ctx, bucket.name, k).Err()
		} else {
			err = db.HSet(ctx, bucket.name, k, new).Err()
		}
		if err == nil {
			for _, f := range endFuncs {
				f()
			}
		}
		return msg, err
	} else {
		return toMaster(bucket.name, k, new)
	}

}

func (bucket *Bucket) GetString(kv ...interface{}) string {
	var key, value string
	for i := range kv {
		if i == 0 {
			key = fmt.Sprint(kv[0])
		} else {
			value = fmt.Sprint(kv[1])
		}
	}
	cs := db.HGet(ctx, bucket.name, key)
	if cs != nil {
		v, _ := cs.Result()
		if v != "" {
			return v
		}
	}
	return value
}

func (bucket *Bucket) GetBytes(key string) (v []byte) {
	cs := db.HGet(ctx, bucket.name, key)
	if cs != nil {
		v, _ = cs.Bytes()
		if len(v) != 0 {
			return
		}
	}
	return
}

func (bucket Bucket) GetInt(key string, vs ...int) int {
	var value int
	if len(vs) != 0 {
		value = vs[0]
	}
	cs := db.HGet(ctx, bucket.name, key)
	if cs != nil {
		v, _ := cs.Result()
		if v != "" {
			return utils.Int(v)
		}
	}
	return value
}

func (bucket Bucket) GetBool(key string, vs ...bool) bool {
	var value bool
	if len(vs) != 0 {
		value = vs[0]
	}
	cs := db.HGet(ctx, bucket.name, key)
	if cs != nil {
		v, _ := cs.Result()
		if v == "true" {
			value = true
		} else if v == "false" {
			value = false
		}
	}
	return value
}

func (bucket Bucket) Foreach(f func(k, v []byte) error) {
	vs, _ := db.HGetAll(ctx, bucket.name).Result()
	for key, value := range vs {
		f([]byte(key), []byte(value))
	}
}
func (bucket Bucket) Create(i interface{}) error {
	sq, err := db.Incr(ctx, bucket.name+"_Sequence").Uint64()
	if err != nil {
		return nil
	}
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	sequence := s.FieldByName("Sequence")
	if _, ok := id.Interface().(int); ok {
		key := id.Int()
		if key == 0 {
			key = int64(sq)
			id.SetInt(key)
		}
		if sequence != reflect.ValueOf(nil) {
			sequence.SetInt(int64(sq))
		}
		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}
		return db.HSet(ctx, bucket.name, key, buf).Err()
	} else {
		key := id.String()
		if key == "" {
			key = fmt.Sprint(sq)
			id.SetString(key)
		}
		if sequence != reflect.ValueOf(nil) {
			sequence.SetInt(int64(sq))
		}
		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}
		return db.HSet(ctx, bucket.name, key, buf).Err()
	}
}

func (bucket Bucket) First(i interface{}) error {
	var err error
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	if v, ok := id.Interface().(int); ok {
		data, _ := db.HGet(ctx, bucket.name, fmt.Sprint(v)).Bytes()
		if len(data) == 0 {
			err = errors.New("record not find")
			return err
		}
		return json.Unmarshal(data, i)
	} else {
		data, _ := db.HGet(ctx, bucket.name, id.Interface().(string)).Bytes()
		if len(data) == 0 {
			err = errors.New("record not find")
			return err
		}
		return json.Unmarshal(data, i)
	}
}

func (bucket Bucket) Size() (size int64, err error) {
	var ks []string
	ks, err = db.HKeys(ctx, bucket.name).Result()
	return int64(len(ks)), err
}

func (bucket Bucket) IsEmpty() (empty bool, err error) {
	var ks []string
	ks, err = db.HKeys(ctx, bucket.name).Result()
	return len(ks) == 0, err
}

func (bucket *Bucket) Keys() ([]string, error) {
	return db.HKeys(ctx, bucket.name).Result()
}
