package user

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
	g "main/app/global"
	"main/app/internal/dao"
	"main/app/internal/model"
	"main/utils/jwt"
	"time"
)

// SUser 是一个结构体，用于处理用户相关的操作
type SUser struct{}

// CheckUserIsExist 检查给定的用户名是否已经存在于数据库中
func (s *SUser) CheckUserIsExist(ctx context.Context, username string) error {
	err := dao.User().User().GetUserByUsername(ctx, username)
	// 如果查找过程中出现错误
	if err != nil {
		// 如果错误不是因为找不到记录
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录错误日志
			g.Logger.Errorf("query [user_subject] record failed, err: %v", err)
			// 返回内部错误
			return fmt.Errorf("internal err")
		}
		// 如果找到了用户名
	} else {
		// 返回用户名已存在的错误
		return fmt.Errorf("username already exist")
	}

	// 如果没有找到记录，返回nil
	return nil
}

// EncryptPassword 对给定的密码进行加密
func (s *SUser) EncryptPassword(password string) string {
	// 使用SHA3算法对密码进行哈希
	d := sha3.Sum224([]byte(password))
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(d[:])
}

// CreateUser 在数据库中创建一个新用户
func (s *SUser) CreateUser(ctx context.Context, userSubject *model.UserSubject) {
	// 在数据库中创建用户
	dao.User().User().CreateUser(ctx, userSubject)
}

// CheckPassword 检查给定的用户名和密码是否匹配
func (s *SUser) CheckPassword(ctx context.Context, userSubject *model.UserSubject) error {
	// 在数据库中查找用户名和密码都匹配的用户
	err := dao.User().User().GetUserByUsernameAndPassword(ctx, userSubject)
	// 如果查找过程中出现错误
	if err != nil {
		// 如果错误不是因为找不到记录
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录错误日志
			g.Logger.Errorf("query [user_subject] record failed, err: %v", err)
			// 返回内部错误
			return fmt.Errorf("internal err")
			// 如果错误是因为找不到记录
		} else {
			// 返回用户名或密码无效的错误
			return fmt.Errorf("invalid username or password")
		}
	}

	// 如果没有错误，返回nil
	return nil
}

// GenerateToken 为给定的用户生成一个JWT令牌
func (s *SUser) GenerateToken(ctx context.Context, userSubject *model.UserSubject) (string, error) {
	// 获取JWT配置
	jwtConfig := g.Config.Auth.Jwt

	// 创建一个新的JWT对象
	j := jwt.NewJWT(&jwt.Config{
		SecretKey:   jwtConfig.SecretKey,
		ExpiresTime: jwtConfig.ExpiresTime,
		BufferTime:  jwtConfig.BufferTime,
		Issuer:      jwtConfig.Issuer})
	// 创建一个新的声明
	claims := j.CreateClaims(&jwt.BaseClaims{
		Id:         userSubject.Id,
		Username:   userSubject.Username,
		CreateTime: userSubject.CreateTime,
		UpdateTime: userSubject.UpdateTime,
	})

	// 生成一个新的令牌
	tokenString, err := j.GenerateToken(&claims)
	// 如果生成令牌过程中出现错误
	if err != nil {
		// 记录错误日志
		g.Logger.Errorf("generate token failed, %v", err)
		// 返回内部错误
		return "", fmt.Errorf("internal err")
	}

	// 将令牌存储在Redis缓存中，并设置一个过期时间
	err = g.Rdb.Set(ctx,
		fmt.Sprintf("jwt_%d", userSubject.Id),
		tokenString,
		time.Duration(jwtConfig.ExpiresTime)*time.Second).Err()
	// 如果存储过程中出现错误
	if err != nil {
		// 记录错误日志
		g.Logger.Errorf("set [jwt] cache failed, %v", err)
		// 返回内部错误
		return "", fmt.Errorf("internal err")
	}

	// 如果没有错误，返回生成的令牌
	return tokenString, nil
}
