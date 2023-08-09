package core

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/cdle/sillyGirl/utils"
)

var temp *PersistentKeyValueStore
var tempPath = filepath.Join(utils.GetDataHome(), "cache.json")

var cache = func(pre string, sec int) interface{} {
	return map[string]interface{}{
		"set": func(key string, value interface{}, num int) {
			if sec != 0 && num == 0 {
				num = sec
			}
			temp.Set(pre+"_"+key, value, num)
		},
		"get": func(key string, def interface{}) interface{} {
			v := temp.Get(pre + "_" + key)
			if v == nil {
				v = def
			}
			return v
		},
		"delete": func(key string) {
			temp.Delete(pre + "_" + key)
		},
	}
}

type PersistentKeyValueStore struct {
	sync.RWMutex
	data map[string]Value
}

type Value struct {
	Data      interface{} `json:"data"`
	ExpiredAt time.Time   `json:"expired_at"`
}

func NewPersistentKeyValueStore() *PersistentKeyValueStore {
	return &PersistentKeyValueStore{
		data: make(map[string]Value),
	}
}

func (s *PersistentKeyValueStore) Set(key string, value interface{}, dur int) error {
	s.Lock()
	defer s.Unlock()
	if dur == 0 {
		dur = 86400
	}
	expiredAt := time.Now().Add(time.Duration(dur) * time.Second)
	s.data[key] = Value{
		Data:      value,
		ExpiredAt: expiredAt,
	}
	go func() {
		s.RLock()
		defer s.RUnlock()
		defer func() {
			recover()
		}()
		jsonBytes, err := json.Marshal(s.data)
		if err == nil {
			ioutil.WriteFile(tempPath, jsonBytes, 0644)
		}
	}()
	return nil
}

func (s *PersistentKeyValueStore) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.data[key]
	if !ok || time.Now().After(value.ExpiredAt) {
		return nil
	}
	return value.Data
}

func (s *PersistentKeyValueStore) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
	// Serialize data to JSON and write to file
	jsonBytes, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tempPath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *PersistentKeyValueStore) LoadFromFile() error {
	s.Lock()
	defer s.Unlock()

	// Read data from file
	jsonBytes, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return err
	}
	// Deserialize JSON data
	err = json.Unmarshal(jsonBytes, &s.data)
	if err != nil {
		return err
	}
	// Delete expired data
	now := time.Now()
	for key, value := range s.data {
		if now.After(value.ExpiredAt) {
			delete(s.data, key)
		}
	}
	// Serialize data to JSON and write to file
	jsonBytes, err = json.Marshal(s.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tempPath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	temp = NewPersistentKeyValueStore()
	// Load data from file
	temp.LoadFromFile()
	// if err != nil {
	// 	fmt.Println("Failed to load data from file:", err)
	// }
	// // Set a key-value pair
	// err = store.Set("foo", map[string]interface{}{
	// 	"bar": "baz",
	// }, 1000)
	// if err != nil {
	// 	fmt.Println("Failed to set key-value pair:", err)
	// }
	// // Get a value
	// value := store.Get("foo")
	// fmt.Println(value)

	// // Delete a key-value pair
	// err = store.Delete("foo")
	// if err != nil {
	// 	fmt.Println("Failed to delete key-value pair:", err)
	// }
}
