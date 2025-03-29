package global

import (
	"github.com/ligaolin/gin_lin/config"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	Db     *gorm.DB
)
