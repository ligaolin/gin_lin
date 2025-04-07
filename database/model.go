package database

import (
	"errors"
	"fmt"

	"github.com/ligaolin/gin_lin/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	Db *gorm.DB
}

// 创建mysql连接
func NewMysql(cfg MysqlConfig) (*Mysql, error) {
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
			cfg.Charset,
			cfg.ParseTime,
			cfg.Loc,
		)), &gorm.Config{})
	return &Mysql{Db: db}, err
}

func (d *Mysql) Model(id uint, param any, m any) error {
	if id != 0 {
		if err := d.Db.First(m, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("不存在的数据")
			} else {
				return err
			}
		}
	}
	utils.AssignFields(param, m)
	return nil
}

// 添加或编辑
func (d *Mysql) Edit(id uint, id_name string, m any, sa []Same) error {
	err := d.Same(id, id_name, m, sa)
	if err != nil {
		return err
	}
	if id == 0 {
		return d.Db.Create(m).Error
	} else {
		return d.Db.Save(m).Error
	}
}

type Same struct {
	Sql     string
	Message string
}

// 判断数据是否重复
func (d *Mysql) Same(id uint, id_name string, model any, data []Same) error {
	var count int64
	for _, v := range data {
		if id != 0 {
			v.Sql += fmt.Sprintf(" AND %s != %d", id_name, id)
		}
		d.Db.Model(&model).Where(v.Sql).Count(&count)
		if count > 0 {
			return errors.New(v.Message)
		}
	}
	return nil
}

type Update struct {
	ID     uint
	IDName string
	Field  string
	Value  any
}

// 更新
func (d *Mysql) Update(param Update, m any, has []string) (err error) {
	if param.ID == 0 {
		return errors.New("id必须")
	}
	if param.IDName == "" {
		param.IDName = "id"
	}
	if param.Field == "" {
		return errors.New("field必须")
	}
	if !utils.Contains(has, param.Field) {
		return errors.New("field数据不合法")
	}
	if err := d.Db.First(m, param.IDName+" = ?", param.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("不存在的数据")
		} else {
			return err
		}
	}
	return d.Db.Model(m).Where(param.IDName+" = ?", param.ID).Update(param.Field, param.Value).Error
}

type Delete struct {
	ID      any
	PIDName string
	IDName  string
}

// 当没有上级时pid_name和id_name都设为""
func (d *Mysql) Delete(param Delete, m any) ([]uint, error) {
	if param.IDName == "" {
		param.IDName = "id"
	}
	ids, err := utils.ToSliceUint(param.ID, ",")
	if err != nil {
		return nil, err
	}
	if param.PIDName != "" {
		if err = d.FindChildrenId(FindChildrenId{IDs: &ids, PIDName: param.PIDName, IDName: param.IDName}, m); err != nil {
			return nil, err
		}
	}

	if err := d.Db.Delete(&m, ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

type FindChildrenId struct {
	IDs     *[]uint
	PIDName string
	IDName  string
}

func (d *Mysql) FindChildrenId(f FindChildrenId, m any) error {
	var cids []uint
	if err := d.Db.Model(&m).Where(f.PIDName+" IN ?", *f.IDs).Pluck(f.IDName, &cids).Error; err != nil {
		return err
	}
	if len(cids) > 0 {
		err := d.FindChildrenId(FindChildrenId{IDs: &cids, PIDName: f.PIDName, IDName: f.IDName}, m)
		if err != nil {
			return err
		}
		*f.IDs = append(*f.IDs, cids...)
	}
	return nil
}

type First struct {
	ID     uint
	Joins  string
	Select string
	IDName string
}

func (d *Mysql) First(f First, m any) error {
	if f.IDName == "" {
		f.IDName = "id"
	}
	var db = d.Db.Where(f.IDName+" = ?", f.ID)
	if f.Joins != "" {
		db.Joins(f.Joins)
	}
	if f.Select != "" {
		db.Select(f.Select)
	}
	if err := db.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("不存在的数据")
		} else {
			return err
		}
	}
	return nil
}

type ListData struct {
	Data      any   `json:"data"`
	Total     int64 `json:"total"`      // 总数量
	TotalPage int64 `json:"total_page"` // 总页数
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
}

type List struct {
	Page     int
	PageSize int
	Order    string
	Where    string
	Joins    string
	Select   string
}

// 列表
func (d *Mysql) List(l List, m any) (ListData, error) {
	var (
		db       = d.Db.Model(m).Where(l.Where)
		total_db = d.Db.Model(m).Where(l.Where)
		data     = ListData{
			Page:      l.Page,
			PageSize:  l.PageSize,
			Data:      m,
			TotalPage: 1,
		}
	)

	if l.Page != 0 {
		if l.PageSize == 0 {
			l.PageSize = 10
		}
		db.Offset((l.Page - 1) * l.PageSize).Limit(l.PageSize)
	}
	if l.Order != "" {
		db.Order(l.Order)
	}
	if l.Joins != "" {
		db.Joins(l.Joins)
		total_db.Joins(l.Joins)
	}
	if l.Select != "" {
		db.Select(l.Select)
	}
	if err := db.Find(data.Data).Error; err != nil {
		return data, err
	}
	if err := total_db.Count(&data.Total).Error; err != nil {
		return data, err
	}
	if l.Page != 0 {
		data.TotalPage = data.Total / int64(l.PageSize)
		if data.Total%int64(l.PageSize) != 0 {
			data.TotalPage++
		}
	}
	return data, nil
}

func (d *Mysql) Code(n int, field string, m any) (string, error) {
	for {
		code := utils.GenerateRandomAlphanumeric(n)
		if err := d.Db.Where(field+" = ?", code).First(m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果记录不存在，说明生成的 code 是唯一的，可以返回
				return code, nil
			} else {
				return "", fmt.Errorf("查询数据库失败: %w", err)
			}
		}
	}
}
