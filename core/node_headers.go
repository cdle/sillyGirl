package core

import (
	"net/http"

	"github.com/dop251/goja"
)

func MakeHeadersObject(vm *goja.Runtime, header http.Header) *goja.Object {
	obj := vm.NewObject()
	obj.Set("get", func(name string) string {
		return header.Get(name)
	})

	obj.Set("has", func(name string) bool {
		return header.Get(name) != ""
	})

	obj.Set("set", func(name, value string) {
		header.Set(name, value)
	})

	obj.Set("append", func(name, value string) {
		header.Add(name, value)
	})

	obj.Set("delete", func(name string) {
		header.Del(name)
	})

	obj.Set("keys", func() []string {
		keys := make([]string, 0, len(header))
		for k := range header {
			keys = append(keys, k)
		}
		return keys
	})

	obj.Set("values", func() []string {
		values := make([]string, 0, len(header))
		for _, v := range header {
			values = append(values, v...)
		}
		return values
	})

	obj.Set("entries", func() [][2]string {
		entries := make([][2]string, 0, len(header))
		for k, v := range header {
			for _, value := range v {
				entries = append(entries, [2]string{k, value})
			}
		}
		return entries
	})

	obj.Set("forEach", func(callback func(value, name string)) {
		for k, v := range header {
			for _, value := range v {
				callback(value, k)
			}
		}
	})

	return obj
}
