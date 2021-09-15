package core

import (
	"github.com/boltdb/bolt"
)

var sillyGirl Bucket = "sillyGirl"
var db *bolt.DB

type Bucket string

func init() {
	var err error
	db, err = bolt.Open(ExecPath+"/sillyGirl.cache", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func (bucket Bucket) Set(key, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			b, _ = tx.CreateBucket([]byte(bucket))
		}
		b.Put([]byte(key), []byte(value))
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
		value = string(b.Get([]byte(key)))
		return nil
	})
	return value
}

func (bucket Bucket) Foreach(f func(k, v []byte) error) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.ForEach(f)
		return nil
	})
}
