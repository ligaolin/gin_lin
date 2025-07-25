package db

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/ligaolin/gin_lin"
	"gorm.io/gorm"
)

type Model struct {
	Db     *gorm.DB
	Error  error
	Pk     int32
	PkName string
	Model  any
}

func NewModel(db *gorm.DB, model any) *Model {
	return &Model{
		Db:     db,
		PkName: "id",
		Model:  model,
	}
}

func (m *Model) SetPkName(pkName string) *Model {
	if m.Error != nil {
		return m
	}
	m.PkName = pkName
	return m
}

func (m *Model) SetPk(pk int32) *Model {
	if m.Error != nil {
		return m
	}
	m.Pk = pk
	if m.Pk != 0 {
		if err := m.Db.First(m.Model, m.Pk).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				m.Error = errors.New("不存在的数据")
				return m
			} else {
				m.Error = err
				return m
			}
		}
	}
	return m
}

func (m *Model) Copy(param any) *Model {
	if m.Error != nil {
		return m
	}

	if err := copier.Copy(m.Model, param); err != nil {
		m.Error = err
	}
	return m
}

type Same struct {
	Db      *gorm.DB
	Message string
}

// 唯一性判断
func (m *Model) NotSame(sames *[]Same) *Model {
	if m.Error != nil {
		return m
	}

	var count int64
	for _, v := range *sames {
		if err := v.Db.Model(m.Model).Where(fmt.Sprintf("%s != ?", m.PkName), m.Pk).Count(&count).Error; err != nil {
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

func (m *Model) Save() *Model {
	if m.Error != nil {
		return m
	}

	if err := m.Db.Save(m.Model).Error; err != nil {
		m.Error = err
		return m
	}
	return m
}

// 更新
func (m *Model) Update(field string, value any, containsas []string) *Model {
	if m.Error != nil {
		return m
	}

	if field == "" {
		m.Error = errors.New("field必须")
		return m
	}
	if !gin_lin.Contains(containsas, field) {
		m.Error = errors.New("field数据不合法")
		return m
	}
	if err := m.Db.Model(m.Model).Where(m.PkName+" = ?", m.Pk).Update(field, value).Error; err != nil {
		m.Error = err
		return m
	}
	return m
}

// 删除
func (m *Model) Delete(id any) *Model {
	if m.Error != nil {
		return m
	}

	if err := m.Db.Delete(m.Model, id).Error; err != nil {
		m.Error = err
		return m
	}
	return m
}

// 生成唯一随机码
func (m *Model) Code(n int, field string) (string, error) {
	for {
		code := gin_lin.GenerateRandomAlphanumeric(n)
		if err := m.Db.Where(field+" = ?", code).First(m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果记录不存在，说明生成的 code 是唯一的，可以返回
				return code, nil
			} else {
				return "", fmt.Errorf("查询数据库失败: %w", err)
			}
		}
	}
}
