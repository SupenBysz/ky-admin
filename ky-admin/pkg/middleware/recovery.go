package middleware

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 返回一个恢复中间件，用于处理请求过程中的panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否已断开
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取堆栈信息
				stack := string(debug.Stack())
				httpRequest := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String())
				headers := c.Request.Header

				// 记录错误日志
				log := logger.GetLogger()
				log.Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.String("request", httpRequest),
					zap.Any("headers", headers),
					zap.String("stack", stack),
				)

				// 断开的连接，直接返回
				if brokenPipe {
					c.Abort()
					return
				}

				// 返回500错误
				api.Failed(c, api.CodeInternalError, "服务器内部错误")
			}
		}()
		c.Next()
	}
}
