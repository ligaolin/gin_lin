package db

import (
	"fmt"
	"reflect"
)

func (m *Model) FindChildrenID(ids *[]int32, pidName string) *Model {
	if m.Error != nil {
		return m
	}

	var cids []int32
	if err := m.Db.Model(m.Model).Where(pidName+" IN ?", *ids).Pluck(m.PkName, &cids).Error; err != nil {
		m.Error = err
		return m
	}
	if len(cids) > 0 {
		if err := m.FindChildrenID(&cids, pidName).Error; err != nil {
			m.Error = err
			return m
		}
		*ids = append(*ids, cids...)
	}
	return m
}

func (m *Model) FindChildren(pid any, pidName string, childrenName string, order string) *Model {
	if m.Error != nil {
		return m
	}

	// 获取反射值
	sliceValue := reflect.ValueOf(m.Model)
	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		m.Error = fmt.Errorf("m must be a pointer to a slice")
		return m
	}

	// 获取切片元素类型
	elemType := sliceValue.Elem().Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 构建查询

	DB := m.Db.Model(m.Model).Where(fmt.Sprintf("%s = ?", m.PkName), pid)
	if order != "" {
		DB = DB.Order(order)
	}

	// 执行查询
	if err := DB.Find(m).Error; err != nil {
		m.Error = err
		return m
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
		if pidName != "" {
			idField = item.FieldByName(pidName)
		} else {
			idField = item.FieldByName("ID")
		}
		if !idField.IsValid() {
			m.Error = fmt.Errorf("model must have an ID field")
			return m
		}
		id := idField.Interface()

		// 获取 Children 字段
		var childrenField reflect.Value
		if childrenName != "" {
			childrenField = item.FieldByName(childrenName)
		} else {
			childrenField = item.FieldByName("Children")
		}
		if !childrenField.IsValid() || !childrenField.CanSet() {
			m.Error = fmt.Errorf("model must have a Children field that can be set")
			return m
		}

		// 创建子节点切片
		childrenSlice := reflect.New(reflect.SliceOf(elemType)).Interface()

		// 递归查询子节点
		if err := m.FindChildren(id, pidName, childrenName, order).Error; err != nil {
			m.Error = err
			return m
		}

		// 设置 Children 字段
		childrenField.Set(reflect.ValueOf(childrenSlice).Elem())
	}

	return m
}
