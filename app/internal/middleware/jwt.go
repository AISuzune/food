package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	g "main/app/global"
	"main/utils/cookie"
	myjwt "main/utils/jwt"
	"net/http"
	"time"
)

// JWTAuthMiddleware jwt中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		var token string // 定义一个变量用于存储token

		// 从全局配置中获取cookie的配置，并创建一个新的cookie写入器
		cookieConfig := g.Config.Auth.Cookie
		cookieWriter := cookie.NewCookieWriter(&cookie.Config{
			Secret: cookieConfig.Secret,
			Ctx:    c,
			Cookie: http.Cookie{
				Path:     "/",
				Domain:   cookieConfig.Domain,
				MaxAge:   cookieConfig.MaxAge,
				Secure:   cookieConfig.Secure,
				HttpOnly: cookieConfig.HttpOnly,
				SameSite: cookieConfig.SameSite,
			},
		})

		// 从cookie中获取"x-token"的值
		ok := cookieWriter.Get("x-token", &token)
		if token == "" || !ok {
			// 如果token为空或者获取失败，返回401 Unauthorized错误
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "not logged in",
				"ok":   false,
			})
			c.Abort() // 终止后续的处理函数
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		// parseToken 解析token包含的信息
		// 从全局配置中获取jwt的配置，并创建一个新的jwt处理器
		jwtConfig := g.Config.Auth.Jwt
		j := myjwt.NewJWT(&myjwt.Config{
			SecretKey: jwtConfig.SecretKey},
		)

		// 解析token
		mc, err := j.ParseToken(token)
		if err != nil {
			// 如果解析失败，返回400 Bad Request错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"ok":   false,
			})
			c.Abort() // 终止后续的处理函数
			return
		}

		// 如果token即将过期，生成一个新的token，并将其保存到cookie和Redis中
		if mc.ExpiresAt.Unix()-time.Now().Unix() < mc.BufferTime {
			// 更新token的过期时间
			mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(g.Config.Auth.Jwt.ExpiresTime) * time.Second))
			// 生成新的token
			newToken, _ := j.GenerateToken(mc)
			// 解析新的token
			newClaims, _ := j.ParseToken(newToken)
			// 将新的token保存到cookie中
			cookieWriter.Set("x-token", newToken)
			// 将新的token保存到Redis中
			err = g.Rdb.Set(c,
				fmt.Sprintf("jwt_%d", newClaims.BaseClaims.Id),
				newToken,
				time.Duration(jwtConfig.ExpiresTime)*time.Second).Err()
			if err != nil {
				// 如果设置Redis失败，返回500 Internal Server Error错误
				g.Logger.Errorf("set [jwt] cache failed, %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "internal err",
					"ok":   false,
				})
				return
			}
		}

		// 将当前请求的用户ID和用户名保存到请求的上下文c上
		c.Set("id", mc.BaseClaims.Id)
		c.Set("username", mc.BaseClaims.Username)
		c.Next() // 后续的处理函数可以通过c.Get("username")来获取当前请求的用户信息
	}
}
