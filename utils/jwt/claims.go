package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"main/utils/cookie"
	"time"
)

// CustomClaims 包含缓冲时间、注册信息和用户信息
type CustomClaims struct {
	BufferTime           int64 // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
	jwt.RegisteredClaims       // token注册信息
	BaseClaims                 // 用户信息
}

// BaseClaims 包含用户ID、用户名、创建时间和更新时间
type BaseClaims struct {
	Id         int64
	Username   string
	CreateTime time.Time
	UpdateTime time.Time
}

// GetClaims 函数从cookie中获取并解析JWT
func GetClaims(secret string, cookie *cookie.Cookie) (*CustomClaims, error) {
	var token string
	ok := cookie.Get("x-token", &token) // 从cookie中获取"x-token"

	if !ok {
		err := errors.New("get token by cookie failed") // 如果获取失败，返回错误
		return nil, err
	}
	j := NewJWT(&Config{SecretKey: secret}) // 创建一个新的JWT处理器
	claims, err := j.ParseToken(token)      // 解析token
	if err != nil {
		err := errors.New("parse token failed") // 如果解析失败，返回错误
		return nil, err
	}
	return claims, nil // 返回解析后的claims
}

// GetUserInfo 函数获取从jwt解析出来的用户信息
func GetUserInfo(secret string, cookie *cookie.Cookie) (*BaseClaims, error) {
	if cl, err := GetClaims(secret, cookie); err != nil {
		return nil, err // 如果获取claims失败，返回错误
	} else {
		return &cl.BaseClaims, nil // 返回用户信息
	}
}

// GetUserID 函数获取从jwt解析出来的用户ID
func GetUserID(secret string, cookie *cookie.Cookie) (int64, error) {
	if cl, err := GetClaims(secret, cookie); err != nil {
		return -1, err // 如果获取claims失败，返回错误
	} else {
		return cl.BaseClaims.Id, nil // 返回用户ID
	}
}
