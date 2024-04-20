package middleware

import (
	"github.com/gin-gonic/gin"
	g "main/app/global"
	"main/app/internal/model/config"
	"net/http"
)

// checkCors 检查当前的源是否在白名单中
func checkCors(currentOrigin string) *config.CORSWhitelist {
	for _, whitelist := range g.Config.Cors.Whitelist {
		// 从配置中迭代CORS头部并匹配
		if currentOrigin == whitelist.AllowOrigin {
			return &whitelist // 如果匹配，返回白名单配置
		}
	}
	return nil // 如果没有匹配，返回nil
}

// Cors 允许所有CORS请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method                                                                                                           // 获取请求方法
		origin := c.Request.Header.Get("Origin")                                                                                             // 获取请求源
		c.Header("Access-Control-Allow-Origin", origin)                                                                                      // 设置允许的源
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-Sign-Id")            // 设置允许的头部
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")                                                            // 设置允许的方法
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type") // 设置暴露的头部
		c.Header("Access-Control-Allow-Credentials", "true")                                                                                 // 设置允许凭证

		// 如果方法是OPTIONS，终止请求并返回无内容状态
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// CorsByRules 基于配置的逻辑处理请求
func CorsByRules() gin.HandlerFunc {
	// 如果模式是允许所有，返回Cors函数
	if g.Config.Cors.Mode == "allow-all" {
		return Cors()
	}
	return func(c *gin.Context) {
		whitelist := checkCors(c.GetHeader("origin")) // 检查源是否在白名单中

		// 如果通过，添加请求头部
		if whitelist != nil {
			c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)     // 设置允许的源
			c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)   // 设置允许的头部
			c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)   // 设置允许的方法
			c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders) // 设置暴露的头部
			if whitelist.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true") // 如果允许凭证，设置允许凭证
			}
		}

		// 如果没有通过，并且模式是strict-whitelist，且请求不是获取健康检查，拒绝请求
		if whitelist == nil && g.Config.Cors.Mode == "strict-whitelist" && !(c.Request.Method == "GET" && c.Request.URL.Path == "/health") {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			// 无论是否通过，如果方法是OPTIONS，终止请求并返回无内容状态
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}

		// 处理请求
		c.Next()
	}
}
