package utils

import "reflect"

// 将结构体A相同字段赋值给结构体B
func AssignFields(a interface{}, b interface{}) {
	reflectA := reflect.ValueOf(a).Elem() // 获取结构体A的反射值
	reflectB := reflect.ValueOf(b).Elem() // 获取结构体B的反射值

	for i := 0; i < reflectA.NumField(); i++ { // 遍历结构体A的字段
		field := reflectA.Field(i) // 获取结构体A的字段的反射值
		if field.CanSet() {        // 确保字段可以被设置（非只读）
			fieldType := reflectA.Type().Field(i).Type                                          // 获取结构体A的字段类型
			fieldName := reflectA.Type().Field(i).Name                                          // 获取字段名称
			if bField := reflectB.FieldByName(fieldName); bField.IsValid() && bField.CanSet() { // 在结构体B中查找同名字段并确认可以赋值
				// bField.Set(field) // 将结构体A的字段赋值给结构体B的对应字段
				if bField.Type() == fieldType { // 确保字段类型相同
					bField.Set(field) // 将结构体A的字段赋值给结构体B的对应字段
				}
			}
		}
	}
}
