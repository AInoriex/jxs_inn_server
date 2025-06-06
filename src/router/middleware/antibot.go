package middleware

import (
	"eshop_server/src/utils/log"
	"github.com/gin-gonic/gin"
)

// IP风控
func IPStrict() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP地址
		ip := c.ClientIP()
		// 获取用户请求地址
		url_path := c.Request.URL.Path
		// 获取请求头字段
		// user_agent := c.Request.Header.Get("User-Agent")
		// referer := c.Request.Header.Get("Referer")
		// origin := c.Request.Header.Get("Origin")
		// cookie := c.Request.Header.Get("Cookie")
		// content_type := c.Request.Header.Get("Content-Type")
		// host := c.Request.Header.Get("Host")
		// connection := c.Request.Header.Get("Connection")
		// accept := c.Request.Header.Get("Accept")
		// accept_encoding := c.Request.Header.Get("Accept-Encoding")
		// accept_language := c.Request.Header.Get("Accept-Language")
		// cache_control := c.Request.Header.Get("Cache-Control")
		// pragma := c.Request.Header.Get("Pragma")
		// upgrade_insecure_requests := c.Request.Header.Get("Upgrade-Insecure-Requests")
		// x_forwarded_for := c.Request.Header.Get("X-Forwarded-For")

		// 打印日志
		log.Infof("IPStrict params ip:%v url_path:%v", ip, url_path)

		// 检查IP地址是否在白名单中
		// if isIPInWhitelist(ip) {
		// 	c.Next()
		// 	return
		// }

		// 检查IP地址是否在黑名单中
		// if isIPInBlacklist(ip) {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		// 	return
		// }

		// 缓存请求地址-IP地址映射关系
		// TODO 定时器从缓存检测`请求地址-IP地址`映射关系请求频率，超过阈值则加入黑名单
		// CacheRequestIP(c.Request.URL.Path, ip)

		c.Next()
	}
}
