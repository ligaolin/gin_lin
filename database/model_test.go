package database

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/file"
	"github.com/ligaolin/gin_lin/utils"
	"gorm.io/gorm"
)

func db() (*Mysql, error) {
	return NewMysql(MysqlConfig{
		User:      "root",
		Password:  "12345678f",
		Host:      "134.175.182.204",
		Port:      3306,
		DBName:    "wp",
		Charset:   "utf8mb4",
		ParseTime: "True",
		Loc:       "Local",
	})
}

func TestEdit(t *testing.T) {
	router := gin.Default()
	router.POST("/api/model_add", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		db.Db.AutoMigrate(&News{})

		var (
			id    uint
			m     News
			param Param
		)
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		err = db.Model(id, param, &m)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%#v", m)
		utils.Success(nil, "", m)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}

type Param struct {
	ID      uint             `json:"id"`
	Title   string           `json:"title"`
	Type    string           `json:"type"`
	Thumb   *file.UploadFile `json:"thumb"`
	Desc    *string          `json:"desc"`
	Content *string          `json:"content"`
}

type News struct {
	gorm.Model
	Title   string           `gorm:"not null;type:varchar(255)"`
	Type    string           `gorm:"not null;type:varchar(20)"`
	Thumb   *file.UploadFile `gorm:"serializer:json;type:longtext"`
	Desc    *string          `gorm:"type:varchar(255)"`
	Content *string          `gorm:"type:text"`

	Sort  int32  `json:"sort" gorm:"default:100"`
	State string `json:"state" gorm:"default:开启;type:varchar(5)"`
}

type UploadFile struct {
	Name      string
	Extension string
	Path      string
	Url       string
	Size      int64
	Type      string
	Mime      string
}
