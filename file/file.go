package file

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type file struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	IsDir     bool   `json:"is_dir"`
	ModTime   string `json:"mod_time"`
	Path      string `json:"path"`
	Url       string `json:"url"`
	Type      string `json:"type"`
}

func FileList(c *gin.Context, path string, name string, page int, page_size int, domain string) (m map[string]any, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	if name != "" {
		var l []os.DirEntry
		for _, v := range files {
			if strings.Contains(v.Name(), name) {
				l = append(l, v)
			}
		}
		files = l
	}
	base, err := FileBase(c, domain)
	if err != nil {
		return
	}
	var (
		start = (page - 1) * page_size
		end   = start + page_size
	)
	if end > len(files)-1 {
		end = len(files)
	}
	if start > len(files)-1 {
		start = len(files) - 1
		if start < 0 {
			start = 0
		}
	}
	var list []file
	for _, v := range files[start:end] {
		var (
			full_name     = v.Name()
			full_name_arr = strings.Split(full_name, ".")
			extension     string
			info          fs.FileInfo
		)

		info, err = v.Info()
		if err != nil {
			return
		}
		full_name_arr = strings.Split(full_name, ".")
		if len(full_name_arr) > 1 {
			extension = full_name_arr[len(full_name_arr)-1]
		}
		mime, _ := FileMimeType(path + "/" + full_name)
		list = append(list, file{
			Name:      full_name_arr[0],
			FullName:  full_name,
			Extension: extension,
			IsDir:     v.IsDir(),
			Size:      info.Size() / 1024,
			ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
			Path:      "/" + path,
			Url:       base + "/" + path + "/" + v.Name(),
			Type:      strings.Split(mime, "/")[0],
		})
	}
	return map[string]any{
		"data":  list,
		"total": len(files),
	}, nil
}

func FileDelete(c *gin.Context, path string) error {
	name := c.QueryArray("name[]")
	for _, v := range name {
		s := strings.ReplaceAll(path+"/"+v, "/..", "")
		if s == "" {
			return fmt.Errorf("路径必须")
		}
		// 删除文件夹及其内容
		err := os.RemoveAll(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建文件目录
func FileMkDir(path string) error {
	// 获取文件所在的目录
	dir := filepath.Dir(path)
	// 创建目录（如果不存在）
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func FileBase(c *gin.Context, domain string) (string, error) {
	if domain == "" {
		ip, port, err := net.SplitHostPort(c.Request.Host)
		if err != nil {
			return "", err
		}
		domain = "http://" + ip + ":" + port
	}
	return domain, nil
}

func FileMimeType(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	mime := http.DetectContentType(buffer)
	return mime, nil
}
