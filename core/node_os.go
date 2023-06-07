package core

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dop251/goja"
)

func getJsOs(vm *goja.Runtime, running func() bool) *goja.Object {
	var jsos = vm.NewObject()
	jsos.Set("readFile", func(name string) []byte {
		data, err := os.ReadFile(name)
		if err != nil {
			panic(Error(vm, err))
		}
		return data
	})
	jsos.Set("readFileSync", func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			panic(Error(vm, err))
		}
		return string(data)
	})
	jsos.Set("writeFileSync", func(path, content string) {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			panic(Error(vm, err))
		}
	})
	jsos.Set("walkFilePath", func(root string, callback func(path string, info os.FileInfo) bool) {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !running() {
				return errors.New("over")
			}
			if err != nil {
				panic(Error(vm, err))
			}
			if !callback(path, info) {
				return errors.New("over")
			}
			return nil
		})
		if err != nil && err.Error() != "over" {
			panic(Error(vm, err))
		}
	})
	jsos.Set("userHomeDir", func() string {
		dir, err := os.UserHomeDir()
		if err != nil {
			panic(Error(vm, err))
		}
		return dir
	})
	jsos.Set("name", runtime.GOOS)
	jsos.Set("arch", runtime.GOARCH)
	return jsos
}
