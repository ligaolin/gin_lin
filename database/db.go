package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

type Mysql struct {
	Db *gorm.DB
}

func NewMysql(cfg MysqlConfig) (*Mysql, error) {
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
			cfg.Charset,
			cfg.ParseTime,
			cfg.Loc,
		)), &gorm.Config{})
	return &Mysql{Db: db}, err
}
