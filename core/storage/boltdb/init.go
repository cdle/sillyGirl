package boltdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cdle/sillyGirl/core/logs"
	"github.com/cdle/sillyGirl/core/storage"

	"github.com/cdle/sillyGirl/utils"
)

type Bucket struct {
	name string
}

var db *bolt.DB

func (b Bucket) String() string {
	return b.name
}

var expirationA = Bucket{
	name: "timeouts",
}

type Expire struct {
	Value string    `json:"value"`
	Time  time.Time `json:"deadline"`
}

func Set(key string, value string, expiration time.Duration) error {
	_, err := expirationA.Set(key, utils.JsonMarshal(&Expire{
		Value: value,
		Time:  time.Now().Add(expiration),
	}))
	return err
}

func Get(key string) string {
	e := Expire{}
	data := expirationA.GetBytes(key)
	json.Unmarshal(data, &e)
	if e.Time.Before(time.Now()) {
		return ""
	}
	return e.Value
}

var Buckets = []Bucket{}

func InitsillyGirl() storage.Bucket {
	bd := utils.GetDataHome() + "sillyGirl.db"
	_, err := os.Stat(bd)
	if err != nil {
		f, err := os.Create(bd)
		if err != nil {
			logs.Info("傻妞无法创建数据文件 %s ，请手动创建", bd)
			os.Exit(0)
		}
		f.Close()
	}
	db, err = bolt.Open(bd, 0600, nil)
	if err != nil {
		logs.Info("傻妞无法创建数据文件 %s ，请手动创建", bd)
		os.Exit(0)
	}

	v := &Bucket{
		name: "sillyGirl",
	}
	return v
}

func (bucket *Bucket) GetName() string {
	return bucket.name
}

func (bucket *Bucket) Copy(name string) storage.Bucket {
	return &Bucket{name: name}
}

func (bucket *Bucket) Type() string {
	return "boltdb"
}

func (bucket *Bucket) Buckets() []string {
	var r []string
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			r = append(r, string(name))
			return nil
		})
		return nil
	})
	return r
}

func (bucket *Bucket) Delete() error {
	err := db.Update(func(tx *bolt.Tx) error {
		k := fmt.Sprint(bucket.name)
		e := tx.DeleteBucket([]byte(k))
		if e == bolt.ErrBucketNotFound {
			return nil
		}
		return e
	})
	return err
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
	return msg, db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket.name))
		if err != nil {
			return err
		}
		if new == "" {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		} else {
			if err := b.Put([]byte(k), []byte(new)); err != nil {
				return err
			}
		}
		return nil
	})
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
	var handles []func(string, string, string) *storage.Final
	for _, listen := range storage.Listens {
		if listen.Name == bucket.name && (listen.Key == key || listen.Key == "*") {
			handles = append(handles, listen.Handle)
		}
	}
	if len(handles) > 0 {
		old := bucket.GetString(key)
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
			}
		}
	}
	return msg, db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket.name))
		if err != nil {
			return err
		}
		if new == "" {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		} else {
			if err := b.Put([]byte(k), []byte(new)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bucket *Bucket) GetString(kv ...interface{}) string {
	// logs.Info(kv)
	var key, value string
	for i := range kv {
		if i == 0 {
			key = fmt.Sprint(kv[0])
		} else {
			value = fmt.Sprint(kv[1])
		}
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b == nil {
			return nil
		}
		if v := string(b.Get([]byte(key))); v != "" {
			value = v
		}
		return nil
	})
	return value
}

func (bucket *Bucket) GetBytes(key string) []byte {
	var value []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b == nil {
			return nil
		}
		if v := b.Get([]byte(key)); v != nil {
			value = v
		}
		return nil
	})
	return value
}

func (bucket *Bucket) GetInt(key string, vs ...int) int {
	var value int
	if len(vs) != 0 {
		value = vs[0]
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b == nil {
			return nil
		}
		v := utils.Int(string(b.Get([]byte(fmt.Sprint(key)))))
		if v != 0 {
			value = v
		}
		return nil
	})
	return value
}

func (bucket *Bucket) GetBool(key string, vs ...bool) bool {
	var value bool
	if len(vs) != 0 {
		value = vs[0]
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b == nil {
			return nil
		}
		v := string(b.Get([]byte(fmt.Sprint(key))))
		if v == "true" {
			value = true
		} else if v == "false" {
			value = false
		}
		return nil
	})
	return value
}

func (bucket *Bucket) Foreach(f func(k, v []byte) error) {
	var bs = [][][]byte{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				bs = append(bs, [][]byte{k, v})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	for i := range bs {
		f(bs[i][0], bs[i][1])
	}
}

func (bucket *Bucket) Keys() ([]string, error) {
	var bs = []string{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b != nil {
			err := b.ForEach(func(k, _ []byte) error {
				if string(k) == "plugins" {
					return nil
				}
				bs = append(bs, string(k))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return bs, err
}

func (bucket *Bucket) Create(i interface{}) error {
	// logs.Info("-", i)
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	sequence := s.FieldByName("Sequence")
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket.name))
		if err != nil {
			return err
		}
		if _, ok := id.Interface().(int); ok {
			key := id.Int()
			sq, err := b.NextSequence()
			if err != nil {
				return err
			}
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
			return b.Put(utils.Itob(uint64(key)), buf)
		} else {
			key := id.String()
			sq, err := b.NextSequence()
			if err != nil {
				return err
			}
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
			return b.Put([]byte(key), buf)
		}
	})
}

func (bucket *Bucket) First(i interface{}) error {
	var err error
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	if v, ok := id.Interface().(int); ok {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket.name))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(fmt.Sprint(v)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			return json.Unmarshal(data, i)
		})
	} else {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket.name))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(id.Interface().(string)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			return json.Unmarshal(data, i)
		})
	}
	return err
}

func (bucket *Bucket) Size() (int64, error) {
	var r int64
	bucket.Foreach(func(k, v []byte) error {
		r++
		return nil
	})
	return r, nil
}

func (bucket *Bucket) IsEmpty() (bool, error) {
	r := true
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket.name))
		if b != nil {
			o, _ := b.Cursor().First()
			r = o == nil
		}
		return nil
	})
	return r, nil
}
