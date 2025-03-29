package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ligaolin/gin_lin/cache"
	"github.com/ligaolin/gin_lin/config"
	"github.com/ligaolin/gin_lin/database"
	"github.com/ligaolin/gin_lin/global"
)

func Run(config_path string) *gin.Engine {
	var err error

	// 初始化配置
	global.Config, err = config.LoadConfig(config_path)
	if err != nil {
		panic(err)
	}

	// 初始化Mysql
	global.Db, err = database.MysqlInit()
	if err != nil {
		panic(err)
	}

	// 清理文件缓存
	go cache.CleanDiskCacheCron()

	// 初始化Gin
	r := gin.Default()
	r.MaxMultipartMemory = global.Config.MaxMultipartMemory << 20

	// 设置静态文件服务
	r.Static("/"+global.Config.Static, "./"+global.Config.Static)

	return r
}
