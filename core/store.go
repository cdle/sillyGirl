package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"

	"github.com/boltdb/bolt"
)

var sillyGirl Bucket = NewBucket("sillyGirl")
var db *bolt.DB

type Bucket string

func (b Bucket) String() string {
	return string(b)
}

var Buckets = []Bucket{}

func NewBucket(name string) Bucket {
	b := Bucket(name)
	Buckets = append(Buckets, b)
	return b
}

func initStore() {
	var err error
	if runtime.GOOS == "darwin" {
		db, err = bolt.Open("./sillyGirl.cache", 0600, nil)
	} else {
		db, err = bolt.Open(dataHome+"/sillyGirl.cache", 0600, nil)
	}
	if err != nil {
		panic(err)
	}
}

func (bucket Bucket) Set(key interface{}, value interface{}) {

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			b, _ = tx.CreateBucket([]byte(bucket))
		}
		k := fmt.Sprint(key)
		v := fmt.Sprint(value)
		if v == "" {
			b.Delete([]byte(k))
		} else {
			b.Put([]byte(k), []byte(v))
		}
		return nil
	})
}

func (bucket Bucket) Get(kv ...interface{}) string {
	var key, value string
	for i := range kv {
		if i == 0 {
			key = fmt.Sprint(kv[0])
		} else {
			value = fmt.Sprint(kv[1])
		}
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
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

func (bucket Bucket) GetInt(key interface{}, vs ...int) int {
	var value int
	if len(vs) != 0 {
		value = vs[0]
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := Int(string(b.Get([]byte(fmt.Sprint(key)))))
		if v != 0 {
			value = v
		}
		return nil
	})
	return value
}

func (bucket Bucket) GetBool(key interface{}, vs ...bool) bool {
	var value bool
	if len(vs) != 0 {
		value = vs[0]
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
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

func (bucket Bucket) Foreach(f func(k, v []byte) error) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			b.ForEach(f)
		}
		return nil
	})
}

var Int = func(s interface{}) int {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return i
}

var Int64 = func(s interface{}) int64 {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return int64(i)
}

func (bucket Bucket) Create(i interface{}) error {
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	sequence := s.FieldByName("Sequence")
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			b, _ = tx.CreateBucket([]byte(bucket))
		}
		if _, ok := id.Interface().(int); ok {
			key := id.Int()
			sq, _ := b.NextSequence()
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
			return b.Put(itob(uint64(key)), buf)
		} else {
			key := id.String()
			sq, _ := b.NextSequence()
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

func itob(i uint64) []byte {
	return []byte(fmt.Sprint(i))
}

func (bucket Bucket) First(i interface{}) error {
	var err error
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	if v, ok := id.Interface().(int); ok {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(fmt.Sprint(v)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			json.Unmarshal(data, i)
			return nil
		})
	} else {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				err = errors.New("bucket not find")
				return nil
			}
			data := b.Get([]byte(id.Interface().(string)))
			if len(data) == 0 {
				err = errors.New("record not find")
				return nil
			}
			json.Unmarshal(data, i)
			return nil
		})
	}

	return err
}

func (bucket Bucket) Find(is []interface{}) error {
	var err error
	// is = append(is, interface{})
	// reflect.ValueOf(is).
	// id := reflect.ValueOf(i).Elem().FieldByName("ID").Int()
	// db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte(bucket))
	// 	if b == nil {
	// 		err = errors.New("bucket not find")
	// 		return nil
	// 	}
	// 	data := b.Get(itob(uint64(id)))
	// 	if len(data) == 0 {
	// 		err = errors.New("record not find")
	// 		return nil
	// 	}
	// 	json.Unmarshal(data, i)
	// 	return nil
	// })
	return err
}

// func (bucket Bucket) Find(o interface{}) {

// }
