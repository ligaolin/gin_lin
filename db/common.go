package db

type MysqlConfig struct {
	User      string `json:"user" toml:"user" yaml:"user"`
	Password  string `json:"password" toml:"password" yaml:"password"`
	Host      string `json:"host" toml:"host" yaml:"host"`
	Port      int    `json:"port" toml:"port" yaml:"port"`
	DBName    string `json:"db_name" toml:"db_name" yaml:"db_name"`
	Charset   string `json:"charset" toml:"charset" yaml:"charset"`
	ParseTime string `json:"parse_time" toml:"parse_time" yaml:"parse_time"`
	Loc       string `json:"loc" toml:"loc" yaml:"loc"`
}

// 更新参数
type UpdateParam struct {
	ID    uint   `json:"id" validate:"required:主键值必须"`
	Field string `json:"field" validate:"required:字段名必须"`
	Value any    `json:"value"`
}

type DeleteParam struct {
	ID any `json:"id" validate:"required:主键值必须"`
}

type FirstParam struct {
	ID uint `form:"id" validate:"required:主键值必须"`
}
type ListBase struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Order    string `form:"order"`
}

// 模型基础字段
type IDCreatedAtUpdatedAt struct {
	ID        uint  `json:"id" gorm:"primarykey"`
	CreatedAt *Time `json:"created_at"`
	UpdatedAt *Time `json:"updated_at"`
}

// 模型基础字段
type IDCreatedAtUpdatedAtDeletedAt struct {
	IDCreatedAtUpdatedAt
	DeletedAt *Time `json:"deleted_at" gorm:"index"`
}

// 模型排序
type Sort struct {
	Sort int32 `json:"sort" gorm:"default:100"`
}

// 模型状态
type State struct {
	State string `json:"state" gorm:"default:开启;type:varchar(5)"`
}

// 模型排序和状态
type SortState struct {
	Sort
	State
}

// 模型基础字段
type IDCreatedAtUpdatedAtDeletedAtSortState struct {
	IDCreatedAtUpdatedAtDeletedAt
	SortState
}

// 模型基础字段
type IDCreatedAtUpdatedAtSortState struct {
	IDCreatedAtUpdatedAt
	SortState
}

type HasChildrenStruct struct {
	HasChildren bool `json:"hasChildren" gorm:"-:all;default:false"`
}

type ChildrenStruct[T any] struct {
	Children []T `json:"children" gorm:"-:all;default:false"`
}
