package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 定义Cookie和Config类型
type (
	Cookie struct {
		Config *Config
	}

	Config struct {
		Secret      string       // 密钥
		Ctx         *gin.Context // Gin上下文
		http.Cookie              // HTTP cookie
	}
)

// NewCookieWriter 函数创建一个新的Cookie写入器
func NewCookieWriter(config *Config) *Cookie {
	return &Cookie{
		Config: config,
	}
}

// Set 方法设置一个新的cookie
func (c *Cookie) Set(key string, value interface{}) {
	bytes, _ := json.Marshal(value)        // 将值转换为JSON
	setSecureCookie(c, key, string(bytes)) // 设置安全cookie
}

// Get 方法获取一个cookie
func (c *Cookie) Get(key string, obj interface{}) bool {
	tempData, ok := getSecureCookie(c, key) // 获取安全cookie
	if !ok {
		return false
	}
	_ = json.Unmarshal([]byte(tempData), obj) // 将cookie值解析为对象
	return true
}

// Remove 方法删除一个cookie
func (c *Cookie) Remove(key string, value interface{}) {
	bytes, _ := json.Marshal(value)        // 将值转换为JSON
	setSecureCookie(c, key, string(bytes)) // 设置安全cookie
}

// setSecureCookie函数设置一个安全的cookie
func setSecureCookie(c *Cookie, name, value string) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))    // 将值进行Base64编码
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10) // 获取当前时间戳
	h := hmac.New(sha256.New, []byte(c.Config.Secret))        // 创建一个新的HMAC哈希
	_, _ = fmt.Fprintf(h, "%s%s", vs, timestamp)              // 将值和时间戳添加到哈希

	sig := fmt.Sprintf("%02x", h.Sum(nil))                    // 计算哈希的签名
	cookie := strings.Join([]string{vs, timestamp, sig}, "|") // 将值、时间戳和签名连接成一个字符串

	// 设置HTTP cookie
	http.SetCookie(c.Config.Ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    cookie,
		MaxAge:   c.Config.MaxAge,
		Path:     "/",
		Domain:   c.Config.Domain,
		SameSite: http.SameSite(1),
		Secure:   c.Config.Secure,
		HttpOnly: c.Config.HttpOnly,
	})
}

// getSecureCookie函数获取一个安全的cookie
func getSecureCookie(c *Cookie, key string) (string, bool) {
	cookie, err := c.Config.Ctx.Request.Cookie(key) // 从请求中获取cookie
	if err != nil {
		return "", false
	}
	val, err := url.QueryUnescape(cookie.Value) // 对cookie值进行解码
	if val == "" || err != nil {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3) // 将cookie值分割成三部分
	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]        // 值
	timestamp := parts[1] // 时间戳
	sig := parts[2]       // 签名

	h := hmac.New(sha256.New, []byte(c.Config.Secret)) // 创建一个新的HMAC哈希
	_, _ = fmt.Fprintf(h, "%s%s", vs, timestamp)       // 将值和时间戳添加到哈希

	// 如果计算的签名和cookie中的签名不匹配，返回false
	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs) // 对值进行Base64解码
	return string(res), true                      // 返回解码后的值
}
