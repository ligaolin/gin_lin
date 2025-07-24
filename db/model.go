package db

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/ligaolin/gin_lin"
	"gorm.io/gorm"
)

func (m *Model) Model(id int32, param any, model any) error {
	if id != 0 {
		if err := m.Db.First(model, id).Error; err != nil {
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
	Same   []Same
}

type Same struct {
	Db      *gorm.DB
	Message string
}

// 唯一性判断
func (m *Model) NotSame(sames *[]Same, id int32, idName string) *Model {
	if m.Error != nil {
		return m
	}
	if idName == "" {
		idName = "id"
	}

	var count int64
	for _, v := range *sames {
		if err := v.Db.Where(fmt.Sprintf("%s = ?", idName), id).Count(&count).Error; err != nil {
			m.Error = err
			return m
		}
		if count > 0 {
			m.Error = errors.New(v.Message)
			return m
		}
	}
	return m
}

type UpdateStruct struct {
	ID     uint
	IDName string
	Field  string
	Value  any
}

// 更新
func (m *Model) Update(param UpdateStruct, model any, containsas []string) *Model {
	if param.ID == 0 {
		return errors.New("id必须")
	}
	if param.IDName == "" {
		param.IDName = "id"
	}
	if param.Field == "" {
		return errors.New("field必须")
	}
	if !gin_lin.Contains(containsas, param.Field) {
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
	ID     any
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
	if l.Select != "" {
		db.Select(l.Select)
	}
	if l.Joins != "" {
		db.Joins(l.Joins)
		total_db.Joins(l.Joins)

		if err := db.Scan(data.Data).Error; err != nil {
			return data, err
		}
	} else {
		if err := db.Find(data.Data).Error; err != nil {
			return data, err
		}
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
