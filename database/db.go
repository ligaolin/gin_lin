package database

import (
	"fmt"

	"github.com/ligaolin/gin_lin/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			global.Config.Mysql.User,
			global.Config.Mysql.Password,
			global.Config.Mysql.Host,
			global.Config.Mysql.Port,
			global.Config.Mysql.DBName,
			global.Config.Mysql.Charset,
			global.Config.Mysql.ParseTime,
			global.Config.Mysql.Loc,
		)), &gorm.Config{})
}
