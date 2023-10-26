package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type FileRouter struct {
}

func NewFileRouter() *FileRouter {
	return &FileRouter{}
}
func checkFile(c *gin.Context, path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		c.JSON(200, gin.H{
			"msg": "文件不存在",
		})
		return file, err
	} else if err != nil {
		c.JSON(200, gin.H{
			"msg": "读取文件失败",
		})
		return nil, err
	}
	return file, nil
}

func (router *FileRouter) Init(r *gin.Engine) {
	defaultRouter := r.Group("/file")
	defaultRouter.GET("/preview", func(c *gin.Context) {
		fmt.Println("com ein")
		path := c.Query("path")
		var err error
		var file *os.File
		path, err = url.QueryUnescape(path)
		file, err = checkFile(c, path)
		if err != nil {
			return
		}
		//判断文件是否存在
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			c.JSON(200, gin.H{
				"msg": "获取文件信息失败",
			})
			return
		}
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		fileName := filepath.Base(path)
		buffer := make([]byte, 512)
		if _, err = file.Read(buffer); err != nil {
			c.JSON(200, gin.H{
				"msg": "获取图片类型失败",
			})
			return
		}
		// 使用http.DetectContentType来自动识别MIME类型
		mimeType := http.DetectContentType(buffer)
		c.Header("Content-Type", mimeType)
		http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
	})

	defaultRouter.GET("/download", func(c *gin.Context) {
		path := c.Query("path")
		var err error
		var file *os.File
		path, err = url.QueryUnescape(path)
		file, err = checkFile(c, path)
		if err != nil {
			return
		}
		defer file.Close()
		//判断文件是否存在
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			c.JSON(200, gin.H{
				"msg": "获取文件信息失败",
			})
			return
		}
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		fileName := filepath.Base(path)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
	})
}
