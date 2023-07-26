package core

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
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

var plugin_dir = utils.ExecPath + "/plugins"
var release = "20230726"

var languages = []Language{
	{
		Name:    "node",
		Version: release,
		Os:      "linux",
		Arch:    "amd64",
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/" + release + "/node_linux_amd64.zip"},
	},
	{
		Name:    "node",
		Version: release,
		Os:      "darwin",
		Arch:    "arm64",
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/" + release + "/node_darwin_arm64.zip"},
	},
}

func initLanguage() {
	// go func() {
	for _, item := range languages {
		if !(item.Os == runtime.GOOS && item.Arch == runtime.GOARCH) {
			continue
		}
		node_dir := utils.ExecPath + "/language/" + item.Name
		os.MkdirAll(node_dir, 0755)
		path := os.Getenv("PATH")
		newPath := ""
		if path != "" {
			newPath = fmt.Sprintf("%s:%s", node_dir, path)
		} else {
			newPath = node_dir
		}
		os.Setenv("PATH", newPath)
		if _, err := os.Stat(node_dir + "/yarn"); err != nil {
			resp, err := http.Get("https://gitee.com/sillybot/binary/releases/download/yarn/yarn.zip")
			if err == nil {
				go func() {
					defer resp.Body.Close()
					zipfile := node_dir + "/yarn.zip"
					f, err := os.OpenFile(zipfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
					if err != nil {
						return
					}
					defer f.Close()
					_, err = io.Copy(f, resp.Body)
					if err != nil {
						return
					}
					defer os.Remove(zipfile)
					unzip(zipfile, 0777, false)
				}()
			}
		}
		func() {
			dir := utils.ExecPath + "/language/" + item.Name
			data, _ := os.ReadFile(dir + "/version")
			if string(data) == item.Version {
				return
			}
			console.Log("正在安装", item.Name, "执行环境....")
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
			if err := unzip(zipfile, 0755, false); err == nil {
				os.WriteFile(dir+"/version", []byte(item.Version), 0755)
			} else {
				// fmt.Println(err)
			}
			console.Log("安装", item.Name, "执行环境成功")
		}()
	}
}

func unzip(filename string, perm fs.FileMode, pkg bool) error {
	zipFile, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	top := ""
	for _, file := range zipFile.File {
		if top == "" {
			top = strings.Split(file.Name, "/")[0]
		}
		// 忽略以 "__MACOSX/" 开头的文件
		if strings.HasPrefix(file.Name, "__MACOSX/") {
			continue
		}
		path := filepath.Join(filepath.Dir(filename), file.Name)
		if file.FileInfo().IsDir() {
			// 如果是目录则创建目录
			err = os.MkdirAll(path, perm)
			if err != nil {
				return err
			}
		} else {
			// 创建文件的父目录
			err = os.MkdirAll(filepath.Dir(path), perm)
			if err != nil {
				return err
			}
			var de = func() error {
				// 创建文件并解压缩数据
				zipFile, err := file.Open()
				if err != nil {
					return err
				}
				defer zipFile.Close()

				localFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
				if err != nil {
					return err
				}
				defer localFile.Close()
				_, err = io.Copy(localFile, zipFile)
				if err != nil {
					return err
				}
				return nil
			}
			if file.Name != top+"/main.js" {
				de()
			} else {
				if pkg {
					defer func() { //安装依赖
						cmd := exec.Command(utils.ExecPath+"/language/node/yarn/bin/yarn", "install")
						cmd.Dir = utils.ExecPath + "/plugins/" + top
						console.Log(cmd.Output())
						de()
					}()
				} else {
					de()
				}

			}
		}
	}
	return nil
}
