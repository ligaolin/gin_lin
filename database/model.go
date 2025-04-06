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

type Query struct {
	Page     *int
	PageSize *int
	Order    *string
	Where    string
	Joins    string
	Select   string
}

func (d *Mysql) List(q Query, model any, data any) (map[string]any, error) {
	var (
		total      int64
		db         = d.Db.Model(model).Where(q.Where)
		total_db   = d.Db.Model(model).Where(q.Where)
		total_page *int64
	)

	if q.Page != nil {
		if q.PageSize == nil {
			var i int = 10
			q.PageSize = &i
		}
		db.Offset((*q.Page - 1) * *q.PageSize).Limit(*q.PageSize)
	}
	if q.Order != nil {
		db.Order(*q.Order)
	}
	if q.Joins != "" {
		db.Joins(q.Joins)
		total_db.Joins(q.Joins)
	}
	if q.Select != "" {
		db.Select(q.Select)
	}

	if err := db.Find(data).Error; err != nil {
		return nil, err
	}
	if err := total_db.Count(&total).Error; err != nil {
		return nil, err
	}
	if q.Page != nil {
		var t = total / int64(*q.PageSize)
		if total%int64(*q.PageSize) != 0 {
			t++
		}
		total_page = &t
	} else {
		q.PageSize = nil
	}
	return map[string]any{
		"data":       data,
		"total":      total,      // 总数量
		"total_page": total_page, // 总页数
		"page":       q.Page,
		"page_size":  q.PageSize,
	}, nil
}

func (d *Mysql) Edit(id uint, m any, sa []Same) error {
	err := same(d.Db, sa, id, m)
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

func same(db *gorm.DB, data []Same, id uint, model any) error {
	var (
		count int64
	)
	for _, v := range data {
		if id != 0 {
			v.Sql += fmt.Sprintf(" AND id != %d", id)
		}
		db.Model(&model).Where(v.Sql).Count(&count)
		if count > 0 {
			return errors.New(v.Message)
		}
	}
	return nil
}

// type UpdateParam struct {
// 	Id    uint    `json:"id"`
// 	Field *string `json:"field"`
// 	Value *any    `json:"value"`
// }

// func (d *Mysql) Update(c *gin.Context, has []string, before any, after any) (param UpdateParam, err error) {
// 	if err = c.Bind(&param); err != nil {
// 		return
// 	}
// 	if param.Id == nil || *param.Id == 0 {
// 		err = errors.New("id必须")
// 		return
// 	}
// 	if param.Field == nil || *param.Field == "" {
// 		err = errors.New("field必须")
// 		return
// 	}
// 	if param.Value == nil {
// 		err = errors.New("value必须")
// 		return
// 	}
// 	if !utils.Contains(has, *param.Field) {
// 		err = errors.New("field数据不合法")
// 		return
// 	}
// 	if err = d.Db.First(before, *param.Id).Error; err != nil {
// 		err = errors.New("不存在的数据")
// 		return
// 	}
// 	err = d.Db.Model(after).Where("id = ?", *param.Id).Update(*param.Field, *param.Value).Error
// 	return
// }

type FindChildrenIdParam struct {
	Ids     *[]uint32
	PidName string
	IdName  string
}

func (d *Mysql) FindChildrenId(f FindChildrenIdParam, m any) error {
	var (
		cids []uint32
	)
	if err := d.Db.Model(&m).Where(f.PidName+" IN ?", *f.Ids).Pluck(f.IdName, &cids).Error; err != nil {
		return err
	}
	if len(cids) > 0 {
		err := d.FindChildrenId(FindChildrenIdParam{Ids: &cids, PidName: f.PidName, IdName: f.IdName}, m)
		if err != nil {
			return err
		}
		*f.Ids = append(*f.Ids, cids...)
	}
	return nil
}

type DeleteParam struct {
	Id      string
	PidName string
	IdName  string
}

// 当没有上级时pid_name和id_name都设为""
func (d *Mysql) Delete(param DeleteParam, m any) ([]uint32, error) {
	if param.IdName == "" {
		param.IdName = "id"
	}
	data, err := utils.StringToSliceUint32(param.Id, ",")
	if err != nil {
		return nil, err
	}
	if param.PidName != "" {
		err = d.FindChildrenId(FindChildrenIdParam{Ids: &data, PidName: param.PidName, IdName: param.IdName}, m)
		if err != nil {
			return nil, err
		}
	}

	if err := d.Db.Delete(&m, data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

type FirstParam struct {
	Id     string
	Joins  string
	Select string
	IdName string
}

func (d *Mysql) First(f FirstParam, m any, model any) error {
	if f.IdName == "" {
		f.IdName = "id"
	}
	var (
		db = d.Db.Model(m).Where(f.IdName+" = ?", f.Id)
	)
	if f.Joins != "" {
		db.Joins(f.Joins)
	}
	if f.Select != "" {
		db.Select(f.Select)
	}
	if err := db.First(&model).Error; err != nil {
		return errors.New("不存在的数据")
	}
	return nil
}

func (d *Mysql) Code(field string, n int, model any) (string, error) {
	var code string
	for {
		code = utils.GenerateRandomAlphanumeric(n)
		if err := d.Db.Where(field+" = ?", code).Select("id").First(model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果记录不存在，说明生成的 code 是唯一的，可以返回
				return code, nil
			}
			// 如果是其他错误，返回错误
			return "", fmt.Errorf("查询数据库失败: %w", err)
		}
		// 如果 model.Id 不为 0，说明 code 已存在，继续生成新的 code
	}
}
