package file

import (
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/utils"
)

type UploadFile struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Path      string `json:"path"`
	Url       string `json:"url"`
	Size      int64  `json:"size"`
	Type      string `json:"type"`
	Mime      string `json:"mime"`
}

func Upload(c *gin.Context, file *multipart.FileHeader, path string, l Limit, static string) (f UploadFile, err error) {
	name := strings.Split(file.Filename, ".")
	mime := file.Header.Get("Content-Type")
	types := strings.Split(mime, "/")[0]

	extension := ""
	if len(name) >= 1 {
		extension = name[1]
	}

	if types != "image" && types != "video" {
		types = "other"
	}

	// 上传限制
	err = limit(extension, types, file.Size, l)
	if err != nil {
		return
	}

	path, err = getPath(path, extension, types, static)
	if err != nil {
		return
	}

	size, err := Save(c, file, path, l.Compress)
	if err != nil {
		return
	}

	base, err := FileBase(c, static)
	if err != nil {
		return
	}
	return UploadFile{
		Name:      name[0],
		Extension: extension,
		Path:      "/" + path,
		Url:       base + "/" + path,
		Size:      size,
		Type:      types,
		Mime:      mime,
	}, nil
}

func getPath(path string, extension string, types string, static string) (string, error) {
	if path == "" {
		path = static + "/upload/" + types + "/" + time.Now().Format("2006-01-02")
	} else {
		path = strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/..", "")
	}
	if !utils.StringPreIs(path, static) {
		return "", errors.New("您上传的路径不符合规范")
	}
	path += "/" + fmt.Sprintf("%d", time.Now().UnixNano()) + "." + extension
	FileMkDir(path)
	return path, nil
}

type Limit struct {
	ImageMaxSize int64
	VideoMaxSize int64
	OtherMaxSize int64
	Extension    string
	Compress     bool
}

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
	s, err := utils.StringToSliceString(l.Extension, ",")
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

func Base64(b64 string) error {
	if b64 == "" {
		return errors.New("缺少图片数据")
	}

	// 解码 Base64 字符串
	decodedData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	fmt.Println(decodedData)

	// 保存为图片文件
	fileName := "output.png" // 你可以根据需要修改文件名和扩展名
	err = os.WriteFile(fileName, decodedData, 0644)
	if err != nil {
		return errors.New("保存图片失败")
	}
	return nil
}
