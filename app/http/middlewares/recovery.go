package middlewares

import (
	"net"
	"net/http/httputil"
	"os"
	"shorturl/pkg/logger"
	"shorturl/pkg/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 读取用户请求数据
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				// 链接中断
				var brokenpipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenpipe = true
						}
					}
				}

				if brokenpipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				// 不是链接中断，属于内部错误，开始进行记录堆栈信息
				logger.Error("recovery from panic",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.Stack("stacktrace"),
				)

				// 记录日志完成， 返回500
				response.Abort500(c)
			}
		}()
		c.Next()
	}
}
