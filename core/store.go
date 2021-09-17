package core

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

var sillyGirl Bucket = NewBucket("sillyGirl")
var db *bolt.DB

type Bucket string

var Buckets = []Bucket{}

func NewBucket(name string) Bucket {
	b := Bucket(name)
	Buckets = append(Buckets, b)
	return b
}

func initStore() {
	var err error
	db, err = bolt.Open(ExecPath+"/sillyGirl.cache", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func (bucket Bucket) Set(key string, value interface{}) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			b, _ = tx.CreateBucket([]byte(bucket))
		}
		b.Put([]byte(key), []byte(fmt.Sprint(value)))
		return nil
	})
}

func (bucket Bucket) Get(kv ...string) string {
	var key, value string
	for i := range kv {
		if i == 0 {
			key = kv[0]
		} else {
			value = kv[1]
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

func (bucket Bucket) GetInt(key string, vs ...int) int {
	var value int
	if len(vs) != 0 {
		value = vs[0]
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := Int(string(b.Get([]byte(key))))
		if v != 0 {
			value = v
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

var Int = func(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
