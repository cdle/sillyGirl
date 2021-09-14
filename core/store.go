package core

import (
	"github.com/boltdb/bolt"
)

var name = "sillyGirl"
var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open(ExecPath+"/"+name+".cache", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			tx.CreateBucket([]byte(name))
		}
		return nil
	})
}

func Set(key string, value string) {
	db.Update(func(tx *bolt.Tx) error {
		tx.Bucket([]byte(name)).Put([]byte(key), []byte(value))
		return nil
	})
}

func Get(key string) string {
	value := ""
	db.View(func(tx *bolt.Tx) error {
		value = string(tx.Bucket([]byte(name)).Get([]byte(key)))
		return nil
	})
	return value
}
