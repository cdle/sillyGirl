package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/cdle/sillyGirl/core/storage/boltdb"
	"github.com/cdle/sillyGirl/core/storage/redis"
	"github.com/cdle/sillyGirl/utils"
	"github.com/goccy/go-json"
)

var bkt storage.Bucket
var HttpPort string
var sillyGirl = MakeBucket("sillyGirl")

var Get = func(key string) string {
	return ""
}
var Set = func(key, value string, expiration time.Duration) error {
	return nil
}

var MakeBucketlocker sync.Mutex

func MakeBucket(name string) storage.Bucket {
	MakeBucketlocker.Lock()
	defer MakeBucketlocker.Unlock()
	if bkt == nil {
		utils.ReadYaml(utils.ExecPath+"/conf/", &Config, "https://raw.githubusercontent.com/cdle/sillyGirl/main/conf/demo_config.yaml")
		utils.SlaveMode = Config.SlaveMode
		HttpPort = Config.HttpPort
		if !Config.EnableRedis {
			bkt = boltdb.InitsillyGirl()
			Get = boltdb.Get
			Set = boltdb.Set
			logs.Info("默认使用boltdb进行数据存储。")
		} else {
			bkt = redis.InitsillyGirl(Config.RedisAddr, Config.RedisPassword)
			Get = redis.Get
			Set = redis.Set
			logs.Info("已使用redis进行数据存储。")
		}
		for _, name := range bkt.Buckets() {
			b := bkt.Copy(name)
			keys, err := b.Keys()
			if len(keys) == 0 && err == nil {
				b.Delete()
			}
		}
	}
	if name == "" {
		name = "sillyGirl"
	}
	if name == "silly" || name == "app" {
		name = "sillyGirl"
	}
	return bkt.Copy(name)
}

func TransformBucketKeyValue(v string) interface{} {
	var result interface{}
	if strings.HasPrefix(v, "f:") {
		result, _ = strconv.ParseFloat(strings.Replace(v, "f:", "", 1), 64)
		return result
	}
	if strings.HasPrefix(v, "d:") {
		result = utils.Int(strings.Replace(v, "d:", "", 1))
		return result
	}
	if strings.HasPrefix(v, "b:") {
		result = strings.Replace(v, "b:", "", 1) == "true"
		return result
	}
	if strings.HasPrefix(v, "o:") {
		json.Unmarshal([]byte(strings.Replace(v, "o:", "", 1)), &result)
		return result
	}
	if v == "" {
		return nil
	}
	return v
}

func GetBucketKeyValue(bucket storage.Bucket, ps ...interface{}) interface{} {
	var key interface{}
	var value interface{}
	if len(ps) == 0 {
		return nil
	}
	if len(ps) > 0 {
		key = ps[0]
	}
	if len(ps) > 1 {
		value = ps[1]
	}
	v := bucket.GetString(key)
	var result = TransformBucketKeyValue(v)
	if result == nil {
		return value
	}
	return result
}

func SetBucketKeyValue(bucket storage.Bucket, key interface{}, value interface{}) (string, error) {
	new := ""
	switch value := value.(type) {
	case int, int64, int32, uint:
		new = fmt.Sprintf("d:%d", value)
	case float32, float64:
		new = fmt.Sprintf("f:%f", value)
	case string, []byte:
		new = fmt.Sprintf("%s", value)
	case bool:
		new = fmt.Sprintf("b:%t", value)
	case nil:
		new = ""
	default:
		new = fmt.Sprintf("o:%s", utils.JsonMarshal(value))
	}
	return bucket.Set(key, new)
}
