package core

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cdle/sillyGirl/utils"
)

type Language struct {
	Name    string   // node
	Version string   // 20230725
	Os      string   // linux
	Arch    string   // amd64
	Links   []string // 下载链接
}

var languages = []Language{
	{
		Name:    "node",
		Version: "20230725",
		Os:      "linux",
		Arch:    "amd64",
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/20230725/node_linux_amd64.zip"},
	},
	{
		Name:    "node",
		Version: "20230725",
		Os:      "darwin",
		Arch:    "arm64",
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/20230725/node_darwin_arm64.zip"},
	},
}

func init() {

	go func() {
		for _, item := range languages {
			if !(item.Os == runtime.GOOS && item.Arch == runtime.GOARCH) {
				continue
			}
			func() {
				dir := utils.ExecPath + "/language/" + item.Name
				data, _ := os.ReadFile(dir + "/version")
				if string(data) == item.Version {
					return
				}
				os.MkdirAll(utils.ExecPath+"/language/"+item.Name, 0755)
				resp, err := http.Get(item.Links[0])
				if err != nil {
					return
				}
				defer resp.Body.Close()
				zipfile := dir + "/" + item.Name + ".zip"
				f, err := os.OpenFile(zipfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					return
				}
				defer f.Close()
				_, err = io.Copy(f, resp.Body)
				if err != nil {
					// fmt.Println(err)
					return
				}
				defer os.Remove(zipfile)
				if err := unzip(zipfile, 0755); err == nil {
					os.WriteFile(dir+"/version", []byte(item.Version), 0755)
				} else {
					// fmt.Println(err)
				}
			}()
		}
	}()
}

func unzip(filename string, perm fs.FileMode) error {
	zipFile, err := zip.OpenReader(filename)
	dir := filepath.Dir(filename)
	fmt.Println(filename, err)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	for _, file := range zipFile.File {
		// 忽略以 "__MACOSX/" 开头的文件
		if strings.HasPrefix(file.Name, "__MACOSX/") {
			continue
		}
		zipFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zipFile.Close()
		localFile, err := os.OpenFile(dir+"/"+file.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			return err
		}
		defer localFile.Close()
		_, err = io.Copy(localFile, zipFile)
		if err != nil {
			return err
		}
	}
	return err
}
