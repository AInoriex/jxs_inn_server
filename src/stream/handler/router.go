package handler

import (
	"eshop_server/src/router/middleware"
	"eshop_server/src/utils/config"
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Cors())

	// 设置路由组
	api := router.Group("/v1/steaming")
	{
		api.POST("/upload_streaming_file", UploadStreamingFile)
		api.GET("/player/:filename", StreamingPlayer)
	}

	log.Infof("初始化路由成功, URL：%s", config.CommonConfig.HttpServer.Addr)
	router.Run(config.CommonConfig.StreamServer.Addr)
}
