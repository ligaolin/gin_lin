package db

import (
	"fmt"
	"testing"
)

func TestBackup(t *testing.T) {
	db, err := NewDbBackup(MysqlConfig{}, "backup.sql")
	if err != nil {
		t.Error(fmt.Errorf("连接数据库失败: %w", err))
	}
	err = db.Backup()
	if err != nil {
		t.Error(fmt.Errorf("备份数据库失败: %w", err))
	}

	// err = db.Restore()
	// if err != nil {
	// 	t.Error(fmt.Errorf("恢复数据库失败: %w", err))
	// }
}
