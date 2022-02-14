// Package middlewares Gin 中间件
package middlewares

import (
	"fmt"
	"shorturl/pkg/auth"
	"shorturl/pkg/config"
	"shorturl/pkg/jwt"
	"shorturl/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
		claims, err := jwt.NewJWT().ParserToken(c)

		// JWT 解析失败，有错误发生
		if err != nil {
			response.Unauthorized(c, fmt.Sprintf("请查看 %v 相关的接口认证文档", config.GetString("app.name")))
			return
		}

		// JWT 解析成功，设置用户信息
		appIDExists := auth.AttemptAppID(claims.AppID) // todo 查询appid是否存在
		if !appIDExists {
			response.Unauthorized(c, "不存在该AppID")
			return
		}

		// 将用户信息存入 gin.context 里
		c.Set("current_appid", claims.AppID)

		c.Next()
	}
}
