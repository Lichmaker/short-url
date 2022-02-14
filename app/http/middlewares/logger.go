package middlewares

import (
	"bytes"
	"shorturl/pkg/helpers"
	"shorturl/pkg/logger"
	"shorturl/pkg/traceid"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type responseBodyWrite struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWrite) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 response 内容
		w := &responseBodyWrite{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		start := time.Now()
		c.Next()

		cost := time.Since(start)
		responStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", helpers.MicrosecondsStr(cost)),
			zap.String("trace-id", traceid.TraceID),
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			requestBody, _ := c.GetRawData()
			logFields = append(logFields, zap.String("Request body", string(requestBody)))

			logFields = append(logFields, zap.String("Response body", w.body.String()))
		}

		if responStatus > 400 && responStatus <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
			logger.Warn("HTTP Warning "+cast.ToString(responStatus), logFields...)
		} else if responStatus >= 500 && responStatus <= 599 {
			// 除了内部错误，记录 error
			logger.Error("HTTP Error "+cast.ToString(responStatus), logFields...)
		} else {
			logger.Debug("HTTP Access Log", logFields...)
		}
	}
}
