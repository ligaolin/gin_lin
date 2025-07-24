package db

import (
	"fmt"
	"reflect"
)

func (m *Model) FindChildrenID(ids *[]int32, pidName string) error {
	var cids []int32
	if err := m.Db.Model(m.Model).Where(pidName+" IN ?", *ids).Pluck(m.PkName, &cids).Error; err != nil {
		return err
	}
	if len(cids) > 0 {
		err := m.FindChildrenID(&cids, pidName)
		if err != nil {
			return err
		}
		*ids = append(*ids, cids...)
	}
	return nil
}

type FindChildrenStruct struct {
	PID          any    // 父节点 ID
	PIDName      string // 父节点 ID 字段名
	Where        string // 查询条件
	Order        string // 排序条件
	IDName       string
	ChildrenName string
}

func (m *Model) FindChildren(param FindChildrenStruct, models any) error {
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
	db := m.Db.Model(models)
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
		if err := m.FindChildren(FindChildrenStruct{
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
