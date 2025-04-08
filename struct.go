package gin_lin

import (
	"fmt"
	"reflect"
)

// 将结构体A相同字段赋值给结构体B
func AssignFields(a any, b any) error {
	// 获取 a 和 b 的反射值
	reflectA := reflect.ValueOf(a)
	reflectB := reflect.ValueOf(b)

	// 如果 a 是指针类型，解引用
	if reflectA.Kind() == reflect.Ptr {
		reflectA = reflectA.Elem()
	}
	// 如果 b 是指针类型，解引用
	if reflectB.Kind() == reflect.Ptr {
		reflectB = reflectB.Elem()
	}

	// 检查 a 和 b 是否是结构体类型
	if reflectA.Kind() != reflect.Struct {
		return fmt.Errorf("a must be a struct or a pointer to a struct")
	}
	if reflectB.Kind() != reflect.Struct {
		return fmt.Errorf("b must be a struct or a pointer to a struct")
	}

	// 缓存结构体 B 的字段信息
	bFieldMap := make(map[string]reflect.Value)
	for i := range reflectB.NumField() {
		field := reflectB.Field(i)
		fieldName := reflectB.Type().Field(i).Name
		bFieldMap[fieldName] = field
	}

	// 遍历结构体 A 的字段
	for i := range reflectA.NumField() {
		fieldA := reflectA.Field(i)
		fieldName := reflectA.Type().Field(i).Name

		// 查找结构体 B 中的同名字段
		if fieldB, ok := bFieldMap[fieldName]; ok {
			// 如果字段类型相同，直接赋值
			if fieldB.Type() == fieldA.Type() && fieldB.CanSet() {
				fieldB.Set(fieldA)
			} else if fieldB.Kind() == reflect.Ptr && fieldA.Kind() == reflect.Ptr {
				// 如果字段是指针类型，解引用后赋值
				if fieldB.IsNil() {
					fieldB.Set(reflect.New(fieldB.Type().Elem()))
				}
				fieldB.Elem().Set(fieldA.Elem())
			} else if fieldB.Kind() == reflect.Ptr && fieldA.Kind() != reflect.Ptr {
				// 如果 B 字段是指针类型，而 A 字段不是，解引用后赋值
				if fieldB.IsNil() {
					fieldB.Set(reflect.New(fieldB.Type().Elem()))
				}
				fieldB.Elem().Set(fieldA)
			} else if fieldB.Kind() != reflect.Ptr && fieldA.Kind() == reflect.Ptr {
				// 如果 B 字段不是指针类型，而 A 字段是，解引用后赋值
				fieldB.Set(fieldA.Elem())
			} else if fieldB.Kind() == reflect.Struct && fieldA.Kind() == reflect.Struct {
				// 处理嵌套结构体字段
				if err := AssignFields(fieldA.Addr().Interface(), fieldB.Addr().Interface()); err != nil {
					return fmt.Errorf("failed to assign nested field %s: %v", fieldName, err)
				}
			}
		}
	}

	return nil
}
