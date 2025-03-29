package file

import (
	"image"
	"io"
	"mime/multipart"
	"os"

	"image/jpeg"
	"image/png"

	"github.com/gin-gonic/gin"
)

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
