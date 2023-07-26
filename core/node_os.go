package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/axgle/mahonia"
	"github.com/cdle/sillyGirl/utils"
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
	type ExecRequest struct {
		Dir     string   `json:"dir"`
		Command []string `json:"command"`
		Env     []string `json:"env"`
		Path    string   `json:"path"`
	}
	jsos.Set("path", utils.ExecPath)
	jsos.Set("exec", func(er ExecRequest) string {
		cmd := exec.Command(er.Command[0], er.Command[1:]...)
		cmd.Dir = er.Dir
		cmd.Env = er.Env
		if er.Path != "" {
			cmd.Path = er.Path
		}
		fmt.Println("==", cmd.Path)
		data, err := cmd.Output()
		if err != nil {
			panic(Error(vm, err))
		}
		return string(data)
	})
	jsos.Set("readFileSync", func(path string, decode string) string {
		// 读取文件内容
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(Error(vm, err))
		}
		if decode == "" || decode == "utf-8" {
			return string(content)
		}
		srcCoder := mahonia.NewDecoder(decode)
		srcResult := srcCoder.ConvertString(string(content))
		return srcResult
		// return ConvertToString(string(content), decode, "utf-8")
		// if decode == "" {
		// 	return string(content)
		// }
		// coder := mahonia.NewEncoder(decode)
		// return coder.ConvertString(string(content))

		// // 从文件内容中自动检测编码
		// encoding, name, _ := charset.DetermineEncoding(content, "")
		// logs.Debug("name", name)

		// // 创建一个将编码转换为 UTF-8 的转换器
		// utf8Reader := transform.NewReader(bytes.NewReader(content), encoding.NewDecoder())

		// // 读取转换后的 UTF-8 数据
		// utf8Content, err := ioutil.ReadAll(utf8Reader)
		// if err != nil {
		// 	panic(Error(vm, err))
		// }
		// return string(utf8Content)
	})

	jsos.Set("writeFileSync", func(path, content string, encode string) {
		converted := ""
		if encode == "" || encode == "utf-8" {
			converted = content
		} else {
			tagCoder := mahonia.NewEncoder(encode)
			converted = tagCoder.ConvertString(content)
		}
		// 覆盖写入文件
		err := os.WriteFile(path, []byte(converted), 0644)
		if err != nil {
			panic(Error(vm, err))
		}
		// if encode == "" {
		// 	encode = "utf-8"
		// }
		// if encode == "utf-8" {
		// 	// 覆盖写入文件
		// 	err := os.WriteFile(path, []byte(content), 0644)
		// 	if err != nil {
		// 		panic(Error(vm, err))
		// 	}
		// 	return
		// }
		// // 将UTF-8编码的文本转换为指定编码
		// index := ianaindex.MIME
		// enc, err := index.Encoding(encode)
		// if err != nil {
		// 	panic(Error(vm, err))
		// }
		// converted, _, err := transform.String(enc.NewEncoder(), content)
		// if err != nil {
		// 	panic(Error(vm, err))
		// }
		// // 覆盖写入文件
		// err = os.WriteFile(path, []byte(converted), 0644)
		// if err != nil {
		// 	panic(Error(vm, err))
		// }
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

func ConvertToString(src string, srcCode string, tagCode string) string {
	if srcCode == tagCode {
		return src
	}
	srcResult := ""
	if srcCode == "utf-8" {
		srcResult = src
	} else {
		srcCoder := mahonia.NewDecoder(srcCode)
		srcResult = srcCoder.ConvertString(src)
	}
	if tagCode == "utf-8" {
		return srcResult
	}
	tagCoder := mahonia.NewEncoder(tagCode)
	result := tagCoder.ConvertString(srcResult)
	return result
}
