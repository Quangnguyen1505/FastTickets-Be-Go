package initialize

import (
	"github.com/ntquang/ecommerce/global"
	"github.com/ntquang/ecommerce/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
