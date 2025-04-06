package database

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ligaolin/gin_lin/utils"
)

type Where struct {
	Name     string
	Op       string
	Value    any
	Nullable bool
}

// ToWhere 将 Where 结构体转换为 SQL where 子句
func ToWhere(data []Where) (string, error) {
	var sql []string
	for _, v := range data {
		if v.Nullable || (!v.Nullable && !isNilOrEmpty(v.Value)) {
			// 处理 Value 值
			value, err := formatValue(v.Value)
			if err != nil {
				return "", fmt.Errorf("格式化字段值失败 %s: %w", v.Name, err)
			}

			switch v.Op {
			case "in":
				sql = append(sql, fmt.Sprintf("%s in ('%s')", v.Name, utils.StringToString(value, ",", "','")))
			case "like":
				sql = append(sql, fmt.Sprintf("%s like '%%%s%%'", v.Name, value))
			case "notLike":
				sql = append(sql, fmt.Sprintf("%s not like '%%%s%%'", v.Name, value))
			case "null":
				sql = append(sql, fmt.Sprintf("%s is null", v.Name))
			case "notNull":
				sql = append(sql, fmt.Sprintf("%s is not null", v.Name))
			case "set":
				sql = append(sql, fmt.Sprintf("FIND_IN_SET('%s','%s')", value, v.Name))
			case "!=":
				sql = append(sql, fmt.Sprintf("%s != '%s'", v.Name, value))
			case ">":
				sql = append(sql, fmt.Sprintf("%s > '%s'", v.Name, value))
			case ">=":
				sql = append(sql, fmt.Sprintf("%s >= '%s'", v.Name, value))
			case "<":
				sql = append(sql, fmt.Sprintf("%s < '%s'", v.Name, value))
			case "<=":
				sql = append(sql, fmt.Sprintf("%s <= '%s'", v.Name, value))
			default:
				sql = append(sql, fmt.Sprintf("%s = '%s'", v.Name, value))
			}

		}
	}
	return strings.Join(sql, " AND "), nil
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
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// formatValue 格式化 Value 为字符串
func formatValue(value any) (string, error) {
	if value == nil {
		return "", nil
	}
	v := reflect.ValueOf(value)
	// 如果是指针，解引用
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil
		}
		v = v.Elem()
		value = v.Interface()
	}
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	case []string:
		return "'" + strings.Join(v, "','") + "'", nil
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64:
		return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]"), nil
	default:
		return "", fmt.Errorf("不支持的值类型: %T", v)
	}
}
