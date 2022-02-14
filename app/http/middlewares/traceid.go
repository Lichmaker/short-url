package middlewares

import (
	"shorturl/pkg/traceid"

	"github.com/gin-gonic/gin"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerTraceId := c.GetHeader("Traceid")
		traceid.Boot(headerTraceId)
		c.Next()
	}
}
