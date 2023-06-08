package core

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dop251/goja"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
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
		// 读取文件内容
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(Error(vm, err))
		}

		// 从文件内容中自动检测编码
		encoding, _, _ := charset.DetermineEncoding(content, "")

		// 创建一个将编码转换为 UTF-8 的转换器
		utf8Reader := transform.NewReader(bytes.NewReader(content), encoding.NewDecoder())

		// 读取转换后的 UTF-8 数据
		utf8Content, err := ioutil.ReadAll(utf8Reader)
		if err != nil {
			panic(Error(vm, err))
		}
		return string(utf8Content)
	})

	jsos.Set("writeFileSync", func(path, content string, encode string) {

		encode = strings.ToLower(encode)
		if encode == "" {
			encode = "utf-8"
		}
		// 将UTF-8编码的文本转换为指定编码
		index := ianaindex.MIME
		enc, err := index.Encoding(encode)
		if err != nil {
			panic(Error(vm, err))
		}
		converted, _, err := transform.String(enc.NewEncoder(), content)
		if err != nil {
			panic(Error(vm, err))
		}
		// 覆盖写入文件
		err = os.WriteFile(path, []byte(converted), 0644)
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
