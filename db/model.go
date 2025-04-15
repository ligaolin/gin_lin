package db

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/ligaolin/gin_lin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	Db *gorm.DB
}

// 创建mysql连接
func NewMysql(cfg *MysqlConfig) (*Mysql, error) {
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
	copier.Copy(m, param)
	return nil
}

type EditStruct struct {
	ID     uint
	IDName string
	Same   []SameStruct
}

type SameStruct struct {
	Sql     string
	Message string
}

// 添加或编辑
func (d *Mysql) Edit(param EditStruct, m any) error {
	var count int64
	if param.IDName == "" {
		param.IDName = "id"
	}
	for _, v := range param.Same {
		if param.ID != 0 {
			v.Sql += fmt.Sprintf(" AND %s != %d", param.IDName, param.ID)
		}
		if err := d.Db.Model(m).Where(v.Sql).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New(v.Message)
		}
	}
	return d.Db.Save(m).Error
}

type UpdateStruct struct {
	ID     uint
	IDName string
	Field  string
	Value  any
}

// 更新
func (d *Mysql) Update(param UpdateStruct, m any, has []string) (err error) {
	if param.ID == 0 {
		return errors.New("id必须")
	}
	if param.IDName == "" {
		param.IDName = "id"
	}
	if param.Field == "" {
		return errors.New("field必须")
	}
	if !gin_lin.Contains(has, param.Field) {
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

type DeleteStruct struct {
	ID      any
	PIDName string
	IDName  string
}

// 当没有上级时pid_name和id_name都设为""
func (d *Mysql) Delete(param DeleteStruct, m any) ([]uint, error) {
	if param.IDName == "" {
		param.IDName = "id"
	}
	ids, err := gin_lin.ToSliceUint(param.ID, ",")
	if err != nil {
		return nil, err
	}
	if param.PIDName != "" {
		if err = d.FindChildrenID(FindChildrenIDStruct{IDs: &ids, PIDName: param.PIDName, IDName: param.IDName}, m); err != nil {
			return nil, err
		}
	}

	if err := d.Db.Delete(&m, ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

type FindChildrenIDStruct struct {
	IDs     *[]uint
	PIDName string
	IDName  string
}

func (d *Mysql) FindChildrenID(f FindChildrenIDStruct, m any) error {
	var cids []uint
	if err := d.Db.Model(&m).Where(f.PIDName+" IN ?", *f.IDs).Pluck(f.IDName, &cids).Error; err != nil {
		return err
	}
	if len(cids) > 0 {
		err := d.FindChildrenID(FindChildrenIDStruct{IDs: &cids, PIDName: f.PIDName, IDName: f.IDName}, m)
		if err != nil {
			return err
		}
		*f.IDs = append(*f.IDs, cids...)
	}
	return nil
}

type FirstStruct struct {
	ID     uint
	Joins  string
	Select string
	IDName string
}

func (d *Mysql) First(f FirstStruct, m any) error {
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

type ListStruct struct {
	Page            int
	PageSize        int
	Order           string
	Where           string
	Joins           string
	Select          string
	PIDName         string
	HasChildrenName string
	IDName          string
}

// 列表
func (d *Mysql) List(l ListStruct, m any) (ListData, error) {
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

	if l.PIDName != "" {
		// 使用反射动态处理 HasChildren 字段
		sliceValue := reflect.ValueOf(data.Data).Elem()
		if sliceValue.Kind() == reflect.Slice {
			for i := range sliceValue.Len() {
				item := sliceValue.Index(i)
				if item.Kind() == reflect.Ptr {
					item = item.Elem()
				}

				// 检查是否存在 HasChildren 字段
				var hasChildrenField reflect.Value
				if l.HasChildrenName != "" {
					hasChildrenField = item.FieldByName(l.HasChildrenName)
				} else {
					hasChildrenField = item.FieldByName("HasChildren")
				}
				if hasChildrenField.IsValid() && hasChildrenField.CanSet() {
					var (
						total   int64
						idField reflect.Value
					)
					if l.IDName != "" {
						idField = item.FieldByName(l.IDName)
					} else {
						idField = item.FieldByName("ID")
					}
					if idField.IsValid() {
						id := idField.Interface()
						if err := d.Db.Model(m).Where(l.PIDName+" = ?", id).Count(&total).Error; err != nil {
							return data, err
						}
						if total > 0 {
							hasChildrenField.SetBool(true)
						} else {
							hasChildrenField.SetBool(false)
						}
					}
				}
			}
		}
	}
	return data, nil
}

func (d *Mysql) Code(n int, field string, m any) (string, error) {
	for {
		code := gin_lin.GenerateRandomAlphanumeric(n)
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

type FindChildrenStruct struct {
	PID          any    // 父节点 ID
	PIDName      string // 父节点 ID 字段名
	Where        string // 查询条件
	Order        string // 排序条件
	IDName       string
	ChildrenName string
}

func (d *Mysql) FindChildren(param FindChildrenStruct, m any) error {
	// 获取反射值
	sliceValue := reflect.ValueOf(m)
	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("m must be a pointer to a slice")
	}

	// 获取切片元素类型
	elemType := sliceValue.Elem().Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 构建查询
	db := d.Db.Model(m)
	if param.Where == "" {
		db = db.Where(fmt.Sprintf("%s = ?", param.PIDName), param.PID)
	} else {
		db = db.Where(param.Where+fmt.Sprintf(" AND %s = ?", param.PIDName), param.PID)
	}
	if param.Order != "" {
		db = db.Order(param.Order)
	}

	// 执行查询
	if err := db.Find(m).Error; err != nil {
		return err
	}

	// 递归查询子节点
	slice := sliceValue.Elem()
	for i := range slice.Len() {
		item := slice.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		// 获取当前节点的 ID
		var idField reflect.Value
		if param.IDName != "" {
			idField = item.FieldByName(param.IDName)
		} else {
			idField = item.FieldByName("ID")
		}
		if !idField.IsValid() {
			return fmt.Errorf("model must have an ID field")
		}
		id := idField.Interface()

		// 获取 Children 字段
		var childrenField reflect.Value
		if param.ChildrenName != "" {
			childrenField = item.FieldByName(param.ChildrenName)
		} else {
			childrenField = item.FieldByName("Children")
		}
		if !childrenField.IsValid() || !childrenField.CanSet() {
			return fmt.Errorf("model must have a Children field that can be set")
		}

		// 创建子节点切片
		childrenSlice := reflect.New(reflect.SliceOf(elemType)).Interface()

		// 递归查询子节点
		if err := d.FindChildren(FindChildrenStruct{
			PID:     id,
			PIDName: param.PIDName,
			Where:   param.Where,
			Order:   param.Order,
		}, childrenSlice); err != nil {
			return err
		}

		// 设置 Children 字段
		childrenField.Set(reflect.ValueOf(childrenSlice).Elem())
	}

	return nil
}
