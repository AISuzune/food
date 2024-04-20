package user

import (
	"github.com/gin-gonic/gin"
	g "main/app/global"
	"main/app/internal/model"
	"main/app/internal/service"
	"main/utils/cookie"
	"net/http"
)

// SignApi 定义一个登录API的结构体
type SignApi struct{}

// Register 注册函数
func (a *SignApi) Register(c *gin.Context) {
	// 从请求中获取用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 如果用户名为空，返回错误
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "username cannot be null",
			"ok":   false,
		})
		return
	}
	// 如果密码为空，返回错误
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "password cannot be null",
			"ok":   false,
		})
		return
	}

	// 检查用户名是否已存在
	err := service.User().User().CheckUserIsExist(c, username)
	if err != nil {
		if err.Error() == "internal err" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
				"ok":   false,
			})
		} else if err.Error() == "username already exist" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"ok":   false,
			})
		}

		return
	}

	// 创建一个用户对象
	userSubject := &model.UserSubject{}

	// 加密密码
	encryptedPassword := service.User().User().EncryptPassword(password)

	// 设置用户名和密码
	userSubject.Username = username
	userSubject.Password = encryptedPassword

	// 在数据库中创建用户
	service.User().User().CreateUser(c, userSubject)

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "register successfully",
		"ok":   true,
	})
}

// Login 登录函数
func (a *SignApi) Login(c *gin.Context) {
	// 从请求中获取用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 如果用户名为空，返回错误
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "username cannot be null",
			"ok":   false,
		})
		return
	}
	// 如果密码为空，返回错误
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "password cannot be null",
			"ok":   false,
		})
		return
	}

	// 创建一个用户对象，并设置用户名和加密后的密码
	userSubject := &model.UserSubject{
		Username: username,
		Password: service.User().User().EncryptPassword(password),
	}

	// 检查用户名和密码是否匹配
	err := service.User().User().CheckPassword(c, userSubject)
	if err != nil {
		switch err.Error() {
		case "internal err":
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
				"ok":   false,
			})
		case "invalid username or password":
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"ok":   false,
			})
		}

		return
	}

	// 为用户生成一个令牌
	tokenString, err := service.User().User().GenerateToken(c, userSubject)
	if err != nil {
		switch err.Error() {
		case "internal err":
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
				"ok":   false,
			})
		}

	}

	// 获取cookie配置
	cookieConfig := g.Config.Auth.Cookie
	// 创建一个新的cookie写入器
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

	// 将令牌写入cookie
	cookieWriter.Set("x-token", tokenString)

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "login successfully",
		"ok":   true,
	})
}
