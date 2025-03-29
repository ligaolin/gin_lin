package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/global"
	"github.com/ligaolin/gin_lin/utils"
	"gorm.io/gorm"
)

type Model interface {
	gorm.Model | any
}
type Base interface {
}

func New[T Model, U Base](id *uint32, param U) (T, error) {
	var m T
	if id != nil {
		if err := global.Db.First(&m, *id).Error; err != nil {
			return m, errors.New("不存在的数据")
		}
	}
	utils.AssignFields(&param, &m)
	return m, nil
}

type Where struct {
	Name     string
	Op       string
	Value    *string
	Nullable bool
}

func ToWhere(data []Where) string {
	var sql []string
	for _, v := range data {
		if v.Nullable || (!v.Nullable && v.Value != nil && *v.Value != "") {
			switch v.Op {
			case "in":
				sql = append(sql, fmt.Sprintf("%s in ('%s')", v.Name, utils.StringToString(*v.Value, ",", "','")))
			case "like":
				sql = append(sql, fmt.Sprintf("%s like '%%%s%%'", v.Name, *v.Value))
			case "notLike":
				sql = append(sql, fmt.Sprintf("%s not like '%%%s%%'", v.Name, *v.Value))
			case "null":
				sql = append(sql, fmt.Sprintf("%s is null", v.Name))
			case "notNull":
				sql = append(sql, fmt.Sprintf("%s is not null", v.Name))
			case "set":
				sql = append(sql, fmt.Sprintf("FIND_IN_SET('%s','%s')", *v.Value, v.Name))
			case "!=":
				sql = append(sql, fmt.Sprintf("%s != '%s'", v.Name, *v.Value))
			case ">":
				sql = append(sql, fmt.Sprintf("%s > '%s'", v.Name, *v.Value))
			case ">=":
				sql = append(sql, fmt.Sprintf("%s >= '%s'", v.Name, *v.Value))
			case "<":
				sql = append(sql, fmt.Sprintf("%s < '%s'", v.Name, *v.Value))
			case "<=":
				sql = append(sql, fmt.Sprintf("%s <= '%s'", v.Name, *v.Value))
			default:
				sql = append(sql, fmt.Sprintf("%s = '%s'", v.Name, *v.Value))
			}

		}
	}
	return strings.Join(sql, " AND ")
}

type Query struct {
	Page     *int
	PageSize *int
	Order    *string
	Where    string
	Joins    string
	Select   string
}

func List[T Model, U Model](q Query) (map[string]interface{}, error) {
	var (
		data       []U
		model      T
		total      int64
		db         = global.Db.Model(&model).Where(q.Where)
		total_db   = global.Db.Model(&model).Where(q.Where)
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

	if err := db.Find(&data).Error; err != nil {
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
	return map[string]interface{}{
		"data":       data,
		"total":      total,      // 总数量
		"total_page": total_page, // 总页数
		"page":       q.Page,
		"page_size":  q.PageSize,
	}, nil
}

func Edit[T Model](id *uint32, m T, same_data []Same) (T, string, error) {
	err := same[T](same_data, id)
	if err != nil {
		return m, "", err
	}
	if id == nil {
		err := global.Db.Create(&m).Error
		return m, "添加成功", err
	} else {
		err := global.Db.Save(&m).Error
		return m, "编辑成功", err
	}
}

type Same struct {
	Sql     string
	Message string
}

func same[T Model](data []Same, id *uint32) error {
	var (
		count int64
		model T
	)
	for _, v := range data {
		if id != nil {
			v.Sql += fmt.Sprintf(" AND id != %d", *id)
		}
		global.Db.Model(&model).Where(v.Sql).Count(&count)
		if count > 0 {
			return errors.New(v.Message)
		}
	}
	return nil
}

type UpdateParam struct {
	Id    *uint32      `json:"id"`
	Field *string      `json:"field"`
	Value *interface{} `json:"value"`
}

func Update[T Model](c *gin.Context, has []string) (param UpdateParam, before T, err error) {
	if err = c.Bind(&param); err != nil {
		return
	}
	if param.Id == nil || *param.Id == 0 {
		err = errors.New("id必须")
		return
	}
	if param.Field == nil || *param.Field == "" {
		err = errors.New("field必须")
		return
	}
	if param.Value == nil {
		err = errors.New("value必须")
		return
	}
	if !utils.Contains(has, *param.Field) {
		err = errors.New("field数据不合法")
		return
	}
	if err = global.Db.First(&before, *param.Id).Error; err != nil {
		err = errors.New("不存在的数据")
		return
	}
	var after T
	err = global.Db.Model(&after).Where("id = ?", *param.Id).Update(*param.Field, *param.Value).Error
	return
}

type FindChildrenIdParam struct {
	Ids     *[]uint32
	PidName string
	IdName  string
}

func FindChildrenId[T Model](f FindChildrenIdParam) error {
	var (
		m    T
		cids []uint32
	)
	if err := global.Db.Model(&m).Where(f.PidName+" IN ?", *f.Ids).Pluck(f.IdName, &cids).Error; err != nil {
		return err
	}
	if len(cids) > 0 {
		err := FindChildrenId[T](FindChildrenIdParam{Ids: &cids, PidName: f.PidName, IdName: f.IdName})
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
func Delete[T Model](d DeleteParam) ([]uint32, error) {
	if d.IdName == "" {
		d.IdName = "id"
	}
	data, err := utils.StringToSliceUint32(d.Id, ",")
	if err != nil {
		return nil, err
	}
	if d.PidName != "" {
		err = FindChildrenId[T](FindChildrenIdParam{Ids: &data, PidName: d.PidName, IdName: d.IdName})
		if err != nil {
			return nil, err
		}
	}
	var m T
	if err := global.Db.Delete(&m, data).Error; err != nil {
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

func First[T Model, U Model](f FirstParam) (U, error) {
	if f.IdName == "" {
		f.IdName = "id"
	}
	var (
		m     T
		db    = global.Db.Model(&m).Where(f.IdName+" = ?", f.Id)
		model U
	)
	if f.Joins != "" {
		db.Joins(f.Joins)
	}
	if f.Select != "" {
		db.Select(f.Select)
	}
	if err := db.First(&model).Error; err != nil {
		return model, errors.New("不存在的数据")
	}
	return model, nil
}

func Code[T Model](field string, n int) (string, error) {
	var model T
	var code string

	for {
		code = utils.GenerateRandomAlphanumeric(n)
		if err := global.Db.Where(field+" = ?", code).Select("id").First(&model).Error; err != nil {
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
