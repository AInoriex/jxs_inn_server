package middleware

import (
	"eshop_server/src/utils/log"
	"fmt"
	"github.com/gin-gonic/gin"
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

		latency := time.Since(start) //执行时间
		statusCode := c.Writer.Status()
		log.Infof("end server | %3d | %10v | %s  %s | %15s |",
			statusCode, formatLatency(latency), method, path, clientIP,
		)
	}
}

// 辅助函数：将耗时转换为友好格式（毫秒/秒/分钟）
func formatLatency(d time.Duration) string {
	switch {
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds()) // 小于1秒显示毫秒
	case d < time.Minute:
		return fmt.Sprintf("%.2fs", d.Seconds()) // 1秒~1分钟显示秒（保留两位小数）
	default:
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60) // 大于1分钟显示分+秒
	}
}
