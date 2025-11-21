package handler

import (
	"eshop_server/src/common/api"
	"github.com/gin-gonic/gin"
)

// 健康检查
func HealthPing(c *gin.Context) {
	api.Success(c, "pong")
}
