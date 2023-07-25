package core

import (
	"io"
	"net/http"
	"os"
	"runtime"

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
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/20230725/node_linux_amd64"},
	},
	{
		Name:    "node",
		Version: "20230725",
		Os:      "darwin",
		Arch:    "arm64",
		Links:   []string{"https://gitee.com/sillybot/binary/releases/download/20230725/node_darwin_arm64"},
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
				f, err := os.OpenFile(dir+"/"+item.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					return
				}
				defer f.Close()
				_, err = io.Copy(f, resp.Body)
				if err != nil {
					return
				}
				os.WriteFile(dir+"/version", []byte(item.Version), 0755)
			}()
		}
	}()
}
