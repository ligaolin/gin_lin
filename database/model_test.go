package database

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/file"
	"github.com/ligaolin/gin_lin/utils"
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

type EditParam struct {
	ID      uint             `json:"id"`
	Title   string           `json:"title" validate:"required:标题必须 len=2,:标题长度不能小于2"`
	Type    string           `json:"type"`
	Thumb   *file.UploadFile `json:"thumb"`
	Desc    *string          `json:"desc"`
	Content *string          `json:"content"`
}

func TestEdit(t *testing.T) {
	router := gin.Default()
	router.POST("/api/model", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		db.Db.AutoMigrate(&News{})

		var (
			m     News
			param EditParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 绑定模型数据
		if err = db.Model(param.ID, param, &m); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 添加或编辑数据
		if err = db.Edit(m.ID, "id", &m, []Same{
			{Sql: fmt.Sprintf("title = '%s'", m.Title), Message: "标题已存在"},
		}); err != nil {
			utils.Error(c, err.Error())
			return
		}
		if param.ID == 0 {
			utils.Success(c, "添加成功", m)
			return
		} else {
			utils.Success(c, "编辑成功", m)
			return
		}
	})
	router.Run()
}

func TestUpdate(t *testing.T) {
	router := gin.Default()
	router.PUT("/api/model", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		var (
			m     News
			param UpdateParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 更新数据
		if err = db.Update(Update{
			ID:    param.ID,
			Field: param.Field,
			Value: param.Value,
		}, &m, []string{"title"}); err != nil {
			utils.Error(c, err.Error())
			return
		}
		utils.Success(c, "更新成功", m)
	})
	router.Run()
}

func TestDelete(t *testing.T) {
	router := gin.Default()
	router.DELETE("/api/model", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		var (
			m     News
			param DeleteParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 更新数据
		ids, err := db.Delete(Delete{
			ID: param.ID,
		}, &m)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		utils.Success(c, "删除成功", ids)
	})
	router.Run()
}

func TestFirst(t *testing.T) {
	router := gin.Default()
	router.GET("/api/model/:id", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		var (
			m     News
			param FirstParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 查询数据
		if err = db.First(First{
			ID: param.ID,
		}, &m); err != nil {
			utils.Error(c, err.Error())
			return
		}
		utils.Success(c, "查询完成", m)
	})
	router.Run()
}

func TestModel(t *testing.T) {
	db, err := db()
	if err != nil {
		t.Error(err)
	}
	var m News
	// 查询数据
	code, err := db.Code(5, "title", &m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(code)
}

type ListParam struct {
	Id    string `form:"id"`
	Title string `form:"title"`
	ListBase
}

func TestList(t *testing.T) {
	router := gin.Default()
	router.GET("/api/model", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		var (
			m     []Region
			param ListParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 查询数据
		where, err := ToWhere([]Where{
			{Name: "id", Op: "like", Value: param.Id},
			{Name: "title", Op: "like", Value: param.Title},
		})
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		data, err := db.List(List{
			Page:     param.Page,
			PageSize: param.PageSize,
			Where:    where,
			Order:    param.Order,
			PIDName:  "pid",
		}, m)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}

		utils.Success(c, "查询完成", data)
	})
	router.Run()
}

type FindChildrenParam struct {
	PID   any    `form:"pid"`
	Title string `form:"title"`
	Order string `form:"order"`
}

func TestFindChildren(t *testing.T) {
	router := gin.Default()
	router.GET("/api/model_children", func(c *gin.Context) {
		db, err := db()
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		var (
			m     []Region
			param FindChildrenParam
		)
		// 绑定参数
		if err := utils.Get(c, &param); err != nil {
			utils.Error(c, err.Error())
			return
		}
		// 查询数据
		where, err := ToWhere([]Where{
			{Name: "title", Op: "like", Value: param.Title},
		})
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
		err = db.FindChildren(FindChildren{
			PID:     param.PID,
			PIDName: "pid",
			Where:   where,
			Order:   param.Order,
		}, &m)
		if err != nil {
			utils.Error(c, err.Error())
			return
		}

		utils.Success(c, "查询完成", m)
	})
	router.Run()
}

type Region struct {
	IDCreatedAtUpdatedAt
	Sort
	Title string `json:"title" gorm:"not null;type:varchar(255)"`
	HasChildrenStruct
	ChildrenStruct[Region]
}

type HasChildrenStruct struct {
	HasChildren bool `json:"hasChildren" gorm:"-:all;default:false"`
}

type ChildrenStruct[T any] struct {
	Children []T `json:"children" gorm:"-:all;default:false"`
}
type News struct {
	IDCreatedAtUpdatedAtDeletedAtSortState
	Title   string           `json:"title" gorm:"not null;type:varchar(255)"`
	Type    string           `json:"type" gorm:"not null;type:varchar(20)"`
	Thumb   *file.UploadFile `json:"thumb" gorm:"serializer:json;type:longtext"`
	Desc    *string          `json:"desc" gorm:"type:varchar(255)"`
	Content *string          `json:"content" gorm:"type:text"`
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
