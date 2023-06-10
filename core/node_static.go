package core

import (
	"encoding/base64"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var statics sync.Map

func addStatic(uuid, path string) {
	statics.Store(uuid, path)
}

func remStatic(uuid string) {
	statics.Delete(uuid)
}

func FindFile(c *gin.Context) {
	// 获取文件名
	filename := c.Param("filename")

	statics.Range(func(_, value any) bool {
		path := value.(string)
		// 拼接文件路径
		filepath := strings.ReplaceAll(filepath.Join(path, filename), "\\/", "\\")
		// 判断文件是否存在
		_, err := os.Stat(filepath)
		if err == nil {
			// 文件存在，读取文件内容并返回
			file, err := ioutil.ReadFile(filepath)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return true
			}
			// 根据文件类型设置Content-Type
			contentType := http.DetectContentType(file)
			c.Header("Content-Type", contentType)

			// 返回文件内容
			c.Data(http.StatusOK, contentType, file)
			return false
		} else {
			console.Log(err)
		}
		return true
	})
	// 如果文件不存在，返回404错误
	c.AbortWithStatus(http.StatusNotFound)
}

// Server.GET("/api/file/:filename", FindFile)
// 	Server.GET("/api/decode/:random", Base642Binary)

func Base642Binary(c *gin.Context) {
	random := c.Param("random")
	s, ok := temp.Get("base64_" + random).(string)
	if !ok {
		c.String(http.StatusBadRequest, "Invalid input")
		return
	}
	input := strings.TrimPrefix(s, "base64://")
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid input")
		return
	}
	// 解析图片格式
	_, format, err := image.DecodeConfig(strings.NewReader(string(data)))
	fmt.Println(format, err)
	if err != nil {
		c.Header("Content-Type", "application/octet-stream")
	} else {
		// 根据图片格式设置响应头
		switch format {
		case "jpeg":
			c.Header("Content-Type", "image/jpeg")
		case "png":
			c.Header("Content-Type", "image/png")
		default:
			c.Header("Content-Type", "application/octet-stream")
			return
		}
	}
	c.Data(http.StatusOK, "application/octet-stream", data)
}
