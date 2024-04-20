package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	Config *Config
}

type Config struct {
	SecretKey   string // 密钥
	ExpiresTime int64  // 过期时间，单位：秒
	BufferTime  int64  // 缓冲时间，缓冲时间内会获得新的token刷新令牌，此时一个用户会存在两个有效令牌，但是前端只留一个，另一个会丢失
	Issuer      string // 签发者
}

// 定义一些常见的错误
var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

// NewJWT 函数创建一个新的JWT处理器
func NewJWT(config *Config) *JWT {
	return &JWT{Config: config}
}

// CreateClaims 方法创建一个新的claims
func (j *JWT) CreateClaims(baseClaims *BaseClaims) CustomClaims {
	claims := CustomClaims{
		BufferTime: j.Config.BufferTime,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),                                  // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.Config.ExpiresTime) * time.Second)), // 过期时间
			Issuer:    j.Config.Issuer,                                                                       // 签发者
		},
		BaseClaims: *baseClaims, // 用户信息
	}
	return claims
}

//func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
//	v, err, _ := g.ConcurrencyControl.Do("JWT_"+oldToken, func() (interface{}, error) {
//		return j.GenerateToken(claims)
//	})
//	return v.(string), err
//}

// GenerateToken 方法生成一个新的token
func (j *JWT) GenerateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims) // 创建一个新的token
	signingKey := []byte(j.Config.SecretKey)                    // 获取签名密钥
	return token.SignedString(signingKey)                       // 返回签名后的token字符串
}

// ParseToken 方法解析JWT
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	signingKey := []byte(j.Config.SecretKey) // 获取签名密钥
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return signingKey, nil // 返回签名密钥
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			// 根据错误类型返回对应的错误
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed // token格式错误
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired // token已过期
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet // token尚未生效
			} else {
				return nil, TokenInvalid // token无效
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil // 如果token有效，返回解析后的claims
		}
		return nil, TokenInvalid // 如果token无效，返回错误
	} else {
		return nil, TokenInvalid // 如果token无效，返回错误
	}
}
