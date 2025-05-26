package main

import (
	"eshop_server/src/stream/handler"
	"eshop_server/src/utils/config"
	// "eshop_server/src/utils/db"
	"eshop_server/src/utils/log"
	"fmt"
)

// 流媒体服务入口
func main() {
	// 初始化配置文件
	err := config.InitConfig()
	if err != nil {
		fmt.Printf("初始化配置文件失败: %s", err.Error())
		return
	}

	// 初始化日志
	log.InitLogger()
	defer log.Sync()
	log.Info("初始化日志成功")

	// TODO 初始化数据库
	// db.InitMysqlAll(config.DbConfig.Mysql.Host, config.DbConfig.Mysql.Db, config.DbConfig.Mysql.MaxCon, db.Con_Main, config.CommonConfig.OpenDbLog)
	// log.Info("初始化Mysql数据库成功")

	// 初始化路由
	handler.InitRouter()
}
