package middleware

import (
	"github.com/gin-gonic/gin"
	"eshop_server/src/utils/log"
	"time"
)

// @Title	自定义日志中间件
// @Description	defer 上报埋点信息（需按固定格式透传）
// @Author  AInoriex  (2025/04/23 17:02)
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		log.Infof("begin server | %s  %s | %15s |", method, path, clientIP)

		// body
		c.Next()

		end := time.Now()
		latency := end.Sub(start) //执行时间
		statusCode := c.Writer.Status()
		log.Infof("end server | %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}
