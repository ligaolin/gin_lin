package db

import (
	"fmt"
	"reflect"
	"strings"
)

type Where struct {
	Name     string
	Op       string
	Value    any
	Nullable bool
}

// 将 Where 结构体转换为安全的 SQL where 子句和参数
func GetWhere(data []Where) (string, []any, error) {
	var sqlParts []string
	var params []any

	for _, v := range data {
		if v.Nullable || (!v.Nullable && !isNilOrEmpty(v.Value)) {
			switch v.Op {
			case "in":
				values, err := toSlice(v.Value)
				if err != nil {
					return "", nil, fmt.Errorf("转换 IN 条件值失败 %s: %w", v.Name, err)
				}
				placeholders := make([]string, len(values))
				for i := range values {
					placeholders[i] = "?"
					params = append(params, values[i])
				}
				sqlParts = append(sqlParts, fmt.Sprintf("%s IN (%s)", v.Name, strings.Join(placeholders, ",")))

			case "like":
				sqlParts = append(sqlParts, fmt.Sprintf("%s LIKE ?", v.Name))
				params = append(params, "%"+toString(v.Value)+"%")

			case "notLike":
				sqlParts = append(sqlParts, fmt.Sprintf("%s NOT LIKE ?", v.Name))
				params = append(params, "%"+toString(v.Value)+"%")

			case "null":
				sqlParts = append(sqlParts, fmt.Sprintf("%s IS NULL", v.Name))

			case "notNull":
				sqlParts = append(sqlParts, fmt.Sprintf("%s IS NOT NULL", v.Name))

			case "set":
				sqlParts = append(sqlParts, fmt.Sprintf("FIND_IN_SET(?, %s)", v.Name))
				params = append(params, toString(v.Value))

			case "!=", ">", ">=", "<", "<=":
				sqlParts = append(sqlParts, fmt.Sprintf("%s %s ?", v.Name, v.Op))
				params = append(params, v.Value)

			default: // "="
				sqlParts = append(sqlParts, fmt.Sprintf("%s = ?", v.Name))
				params = append(params, v.Value)
			}
		}
	}

	return strings.Join(sqlParts, " AND "), params, nil
}

// isNilOrEmpty 判断 Value 是否为 nil 或空值
func isNilOrEmpty(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	// 如果是指针，解引用
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// toString 安全转换为字符串
func toString(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toSlice 将值转换为切片
func toSlice(value any) ([]any, error) {
	if value == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return []any{value}, nil
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}

	return result, nil
}
