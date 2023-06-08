package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	iconv "github.com/djimenez/iconv-go"
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
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(Error(vm, err))
		}
		// 自动转换编码为UTF-8
		utf8Content, err := iconv.ConvertString(string(content), "auto", "utf-8")
		if err != nil {
			panic(Error(vm, err))
		}
		return utf8Content
	})

	jsos.Set("writeFileSync", func(path, content string, encode string) {
		// 将文本内容转换为GBK编码
		var err error
		if !(encode == "" || encode == "utf-8") {
			content, err = iconv.ConvertString(content, "utf-8", encode)
			if err != nil {
				panic(Error(vm, err))
			}
		}
		// 将文本内容以GBK编码写入文件
		err = ioutil.WriteFile(path, []byte(content), 0644)
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
