package core

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
		filepath := filepath.Join(path, filename)
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
		}
		return true
	})
	// 如果文件不存在，返回404错误
	c.AbortWithStatus(http.StatusNotFound)
}
