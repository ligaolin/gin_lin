package database

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/ligaolin/gin_lin/file"
	"github.com/ligaolin/gin_lin/global"
)

func Backup(path string) error {
	var backupSQL strings.Builder

	// 获取所有表名
	var database_name string
	if err := global.Db.Raw("SELECT DATABASE()").Scan(&database_name).Error; err != nil {
		return err
	}
	backupSQL.WriteString("CREATE DATABASE IF NOT EXISTS `" + database_name + "`;\n")
	backupSQL.WriteString("USE `" + database_name + "`;\n\n")

	// 获取所有表名
	var tableNames []string
	if err := global.Db.Raw("SHOW TABLES").Scan(&tableNames).Error; err != nil {
		return err
	}

	for _, v := range tableNames {
		// 备份表结构
		sql, err := backupTableStructure(v)
		if err != nil {
			return err
		}
		backupSQL.WriteString(sql)

		// 备份表数据
		tableDataSQL, err := backupTableData(v)
		if err != nil {
			return err
		}
		backupSQL.WriteString(tableDataSQL)
	}

	// 创建目录（如果不存在）
	if err := file.FileMkDir(path); err != nil {
		return err
	}

	// 将备份内容写入文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(backupSQL.String())
	if err != nil {
		return err
	}

	return nil
}

type TableInfo struct {
	Table       string `gorm:"column:Table"`
	CreateTable string `gorm:"column:Create Table"`
}

func backupTableStructure(tableName string) (string, error) {
	var createTableSQL TableInfo
	err := global.Db.Raw("SHOW CREATE TABLE " + tableName).Scan(&createTableSQL).Error
	if err != nil {
		return "", err
	}
	return "DROP TABLE IF EXISTS `" + tableName + "`;\n" + createTableSQL.CreateTable + ";\n\n", nil
}

func backupTableData(tableName string) (string, error) {
	var (
		dataSQL  strings.Builder
		has_data = false
		i        = 0
	)
	rows, err := global.Db.Raw("SELECT * FROM " + tableName).Rows()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	for rows.Next() {
		// 每 1000 条数据生成一个 INSERT 语句
		if i%1000 == 0 {
			if i > 0 {
				dataSQL.WriteString(";\n")
			}
			dataSQL.WriteString("INSERT INTO `" + tableName + "` (`" + strings.Join(columns, "`, `") + "`) VALUES \n\t(")
		} else {
			dataSQL.WriteString(",\n\t(")
		}
		has_data = true
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return "", err
		}
		for i, value := range values {
			if i > 0 {
				dataSQL.WriteString(", ")
			}
			if value == nil {
				dataSQL.WriteString("NULL")
			} else {
				val := reflect.ValueOf(value)
				switch val.Kind() {
				case reflect.String:
					dataSQL.WriteString("'" + fmt.Sprintf("%s", value) + "'")
				case reflect.Struct:
					// 判断是否是 time.Time 类型
					if t, ok := value.(time.Time); ok {
						// 将时间格式化为 SQL 支持的格式
						dataSQL.WriteString("'" + t.Format("2006-01-02 15:04:05") + "'")
					} else {
						dataSQL.WriteString(fmt.Sprintf("'%v'", value))
					}
				default:
					dataSQL.WriteString(fmt.Sprintf("'%s'", value))
				}
			}
		}
		dataSQL.WriteString(")")
		i++
	}

	if has_data {
		return dataSQL.String() + ";\n\n", nil
	} else {
		return "", nil
	}
}

func Reduction(path string) error {
	// 读取SQL文件
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 将SQL文件内容按分号分割成多个SQL语句
	sqlStatements := strings.Split(string(content), ";")

	// 执行每个SQL语句
	for _, sql := range sqlStatements {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		if err := global.Db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
