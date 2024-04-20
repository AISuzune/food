package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// ZapLogger receive log from gin
func ZapLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             // 记录请求开始的时间
		path := c.Request.URL.Path      // 获取请求的路径
		query := c.Request.URL.RawQuery // 获取请求的查询参数
		c.Next()                        // 调用后续的处理函数

		cost := time.Since(start) // 计算请求的耗时
		// 记录请求日志
		logger.Infow(path,
			zap.Int("status", c.Writer.Status()),                                 // 响应状态码
			zap.String("method", c.Request.Method),                               // 请求方法
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // 请求查询参数
			zap.String("ip", c.ClientIP()),                                       // 客户端IP
			zap.String("user-agent", c.Request.UserAgent()),                      // 用户代理
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 错误信息
			zap.Duration("cost", cost),                                           // 请求耗时
		)
	}
}

// ZapRecovery recover incoming panic
func ZapRecovery(logger *zap.SugaredLogger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a condition that warrants a panic stack trace.
				// 检查是否为断开的连接，因为这不是一个真正需要panic堆栈跟踪的情况
				var brokenPipe bool
				var netErr error
				var ne *net.OpError
				// 使用errors.As函数来判断err是否为*net.OpError类型
				if errors.As(netErr, &ne) {
					var se *os.SyscallError
					// 使用errors.As函数来判断ne.Err是否为*os.SyscallError类型
					if errors.As(ne.Err, &se) {
						// 如果错误信息包含"broken pipe"或"connection reset by peer"，则认为是断开的连接
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false) // 获取HTTP请求的转储
				if brokenPipe {
					// 记录错误日志
					logger.Errorw(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，我们不能向其写入状态
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()            // 终止后续的处理函数
					return
				}

				if stack {
					// 记录错误日志和堆栈跟踪
					logger.Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					// 记录错误日志
					logger.Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError) // 返回500 Internal Server Error状态码
			}
		}()

		c.Next() // 调用后续的处理函数
	}
}
