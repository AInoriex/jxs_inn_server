package middleware

import (
	// "net/http"
	// "fmt"
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 需要接收客户端发送的Origin
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, Origin")
			//// 允许浏览器（客户端）可以解析的头部 （重要）
			//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			////设置缓存时间
			//c.Header("Access-Control-Max-Age", "172800")
		}
		
		// 允许类型校验
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
