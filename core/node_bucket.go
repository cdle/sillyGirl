package core

import (
	"github.com/cdle/sillyGirl/core/storage"
	"github.com/dop251/goja"
)

func MakeBucketObject(vm *goja.Runtime, uuid string, on_start bool, bucket storage.Bucket) *goja.Object {
	obj := vm.NewObject()
	obj.Set("get", func(v ...interface{}) interface{} {
		return GetBucketKeyValue(bucket, v...)
	})
	obj.Set("getAll", func(v ...interface{}) map[string]interface{} {
		var rt = map[string]interface{}{}
		bucket.Foreach(func(b1, b2 []byte) error {
			rt[string(b1)] = TransformBucketKeyValue(string(b2))
			return nil
		})
		return rt
	})
	obj.Set("set", func(key, value interface{}) interface{} {
		msg, err := SetBucketKeyValue(bucket, key, value)
		if err != nil {
			panic(Error(vm, err))
		}
		return msg
	})
	obj.Set("delete", func(key interface{}) error {
		_, err := bucket.Set(key, "")
		return err
	})
	obj.Set("deleteAll", func() error {
		return bucket.Delete()
	})
	obj.Set("keys", func() []string {
		keys, err := bucket.Keys()
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return keys
	})
	obj.Set("len", func() int {
		keys, err := bucket.Keys()
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return len(keys)
	})
	obj.Set("buckets", func() []string {
		return bucket.Buckets()
	})
	obj.Set("_name", func() string {
		return bucket.GetName()
	})
	obj.Set("watch", func(key string, f func(old, new interface{}, key string) *storage.Final) {
		if on_start {
			storage.Watch(bucket, key, func(old, new, key string) *storage.Final {
				// mutex := GetMutex(uuid)
				// mutex.Lock()
				// defer mutex.Unlock()
				return f(TransformBucketKeyValue(old), TransformBucketKeyValue(new), key)
			}, uuid)
		}
	})
	return obj
}

// vm.Set("Bucket", func(name string) interface{} {
// 	return vm.NewProxy(MakeBucketObject(vm, uuid, on_start, MakeBucket(name)), &goja.ProxyTrapConfig{
// 		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
// 			return nil
// 		},
// 		Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
// 			return true
// 		},
// 	})
// })

func JsBucket(vm *goja.Runtime, name string, uuid string, on_start bool) goja.Proxy {
	return vm.NewProxy(MakeBucketObject(vm, uuid, on_start, MakeBucket(name)), &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			obj := target.Get(property)
			if obj != nil {
				return obj
			}
			result := target.Get("get").Export().(func(...interface{}) interface{})(property)
			return vm.ToValue(result)
		},
		Set: func(target *goja.Object, property string, value, receiver goja.Value) (success bool) {
			target.Get("set").Export().(func(interface{}, interface{}) interface{})(
				property, value.Export(),
			)
			return true
		},
	})
}
