package file

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin"
)

type Files struct {
	Context *gin.Context
	Config  *FileConfig
}

func NewFile(c *gin.Context, cfg *FileConfig) *Files {
	return &Files{
		Context: c,
		Config:  cfg,
	}
}

func (f *Files) Upload(file *multipart.FileHeader, dir string, l Limit) (*File, error) {
	extension := filepath.Ext(file.Filename)
	baseName := file.Filename
	if len(extension) > 0 {
		extension = extension[1:]
		baseName = file.Filename[:len(file.Filename)-len(extension)-1]
	}
	mime := file.Header.Get("Content-Type")
	types := strings.Split(mime, "/")[0]

	if types != "image" && types != "video" {
		types = "other"
	}

	// 上传限制
	err := limit(extension, types, file.Size, l)
	if err != nil {
		return nil, err
	}

	// 获取文件保存路径
	path, err := f.GetPath(dir, types)
	if err != nil {
		return nil, err
	}
	path += "/" + fmt.Sprintf("%d", time.Now().UnixNano()) + "." + extension

	// 保存文件
	size, err := Save(f.Context, file, path, l.Compress)
	if err != nil {
		return nil, err
	}

	// 获取文件访问域名
	base, err := Domain(f.Context, f.Config.Domain)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:      baseName,
		Extension: extension,
		FullName:  file.Filename,
		Path:      "/" + path,
		Url:       base + "/" + path,
		Size:      size,
		Type:      types,
		IsDir:     false,
		ModTime:   time.Now().Format("2006-01-02 15:04:05"),
		Mime:      mime,
	}, nil
}

// 获取文件保存路径
func (f *Files) GetPath(dir string, types string) (string, error) {
	if dir == "" {
		// 默认路径
		dir = f.Config.Static + "/upload"
		if types != "" {
			dir += "/" + types + "/" + time.Now().Format("2006-01-02")
		}
	} else {
		// 使用提供的路径，去掉文件夹中包含..的目录
		dir = strings.ReplaceAll(strings.TrimPrefix(dir, "/"), "/..", "")
		// 路径必须在静态目录下
		if !gin_lin.StringPreIs(dir, f.Config.Static) {
			return "", errors.New("您上传的路径不符合规范")
		}
	}
	// 创建文件目录
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return dir, nil
}

type Limit struct {
	ImageMaxSize int64
	VideoMaxSize int64
	OtherMaxSize int64
	Extension    string
	Compress     bool
}

// 上传限制
func limit(extension string, types string, size int64, l Limit) error {
	if types == "image" {
		if l.ImageMaxSize*1024*1024 < size {
			return fmt.Errorf("图片不能超过%dM", l.ImageMaxSize)
		}
	} else if types == "video" {
		if l.VideoMaxSize*1024*1024 < size {
			return fmt.Errorf("视频不能超过%dM", l.VideoMaxSize)
		}
	} else {
		if l.OtherMaxSize*1024*1024 < size {
			return fmt.Errorf("文件不能超过%dM", l.OtherMaxSize)
		}
	}
	s, err := gin_lin.StringToSliceString(l.Extension, ",")
	if err != nil {
		return err
	}
	ok := false
	for _, v1 := range s {
		if strings.EqualFold(extension, v1) {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("%s格式不支持上传", extension)
	}
	return nil
}

type ListParam struct {
	Path     string
	Name     string
	Page     int
	PageSize int
}
type ListRes struct {
	Data  []File `json:"data"`
	Total int64  `json:"total"`
}

func (f *Files) List(param ListParam) (*ListRes, error) {
	// 获取文件保存路径
	var err error
	param.Path, err = f.GetPath(param.Path, "")
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(param.Path)
	if err != nil {
		return nil, err
	}

	// 名称模糊查询
	if param.Name != "" {
		var l []os.DirEntry
		for _, v := range files {
			if strings.Contains(v.Name(), param.Name) {
				l = append(l, v)
			}
		}
		files = l
	}
	base, err := Domain(f.Context, f.Config.Domain)
	if err != nil {
		return nil, err
	}
	var (
		start = (param.Page - 1) * param.PageSize // 开始位置
		end   = start + param.PageSize            // 结束位置
		list  []File
	)
	if end > len(files)-1 {
		end = len(files)
	}
	if start > len(files)-1 {
		start = max(len(files)-1, 0)
	}
	for _, v := range files[start:end] {
		info, err := v.Info()
		if err != nil {
			return nil, err
		}

		name := v.Name()
		extension := filepath.Ext(name)
		baseName := name
		if len(extension) > 0 {
			extension = extension[1:]
			baseName = name[:len(name)-len(extension)-1]
		}

		mime, _ := FileMimeType(param.Path + "/" + name)
		types := strings.Split(mime, "/")[0]
		if types != "image" && types != "video" {
			types = "other"
		}

		path := "/" + param.Path + "/" + name
		list = append(list, File{
			Name:      baseName,
			Extension: extension,
			FullName:  name,
			Path:      path,
			Url:       base + path,
			Size:      info.Size() / 1024,
			Type:      types,
			IsDir:     v.IsDir(),
			ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
			Mime:      mime,
		})
	}
	return &ListRes{Data: list, Total: int64(len(files))}, nil
}

func (f *Files) Delete(path string, name string) error {
	if path == "" {
		// 默认路径
		path = f.Config.Static + "/upload"
	} else {
		// 使用提供的路径，去掉文件夹中包含..的目录
		path = strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/..", "")
		// 路径必须在静态目录下
		if !gin_lin.StringPreIs(path, f.Config.Static) {
			return errors.New("您要删除的路径不符合规范")
		}
	}
	if name != "" {
		path += "/" + name
	}

	// 删除文件夹及其内容
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func FileMimeType(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	mime := http.DetectContentType(buffer)
	return mime, nil
}

func Domain(c *gin.Context, domain string) (string, error) {
	if domain == "" {
		ip, port, err := net.SplitHostPort(c.Request.Host)
		if err != nil {
			return "", err
		}
		return "http://" + ip + ":" + port, nil
	} else {
		return domain, nil
	}
}

// 对于常见图片格式进行压缩，并保存文件
func Save(c *gin.Context, file *multipart.FileHeader, path string, compress bool) (size int64, err error) {
	size = file.Size
	fileReader, err := file.Open()
	if err != nil {
		return
	}
	defer fileReader.Close()

	// 解码图像
	img, kind, err := image.Decode(fileReader)
	if err != nil {
		compress = false
	}

	if compress {
		var outFile *os.File
		outFile, err = os.Create(path)
		if err != nil {
			return
		}
		defer outFile.Close()

		switch kind {
		case "jpeg":
			err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
			if err == nil {
				size, err = outFile.Seek(0, io.SeekEnd) // 获取当前文件指针位置
			} else {
				err = c.SaveUploadedFile(file, path)
			}
		case "png":
			err = png.Encode(outFile, img)
			if err == nil {
				size, err = outFile.Seek(0, io.SeekEnd) // 获取当前文件指针位置
			} else {
				err = c.SaveUploadedFile(file, path)
			}
		default:
			err = c.SaveUploadedFile(file, path)
		}
		if err != nil {
			return
		}
	} else {
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			return
		}
	}

	return
}
