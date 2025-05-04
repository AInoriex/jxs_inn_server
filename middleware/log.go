package middleware

import (
	"github.com/gin-gonic/gin"
	"eshop_server/utils/log"
	"time"
	"net/http"
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT,DELETE, UPDATE")
			////允许跨域设置可以返回其他子段，可以自定义字段
			//c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session")
			//// 允许浏览器（客户端）可以解析的头部 （重要）
			//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			////设置缓存时间
			//c.Header("Access-Control-Max-Age", "172800")
			////允许客户端传递校验信息比如 cookie (重要)
			//c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

