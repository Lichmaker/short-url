package bootstrap

import (
	"net/http"
	"shorturl/app/http/middlewares"
	"shorturl/routes"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRoute(router *gin.Engine) {
	// 中间件
	registerGlobalMiddleWare(router)

	// 404页面
	setup404Handler(router)

	// 配置api路由
	routes.RegisterAPIRoutes(router)
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		middlewares.TraceID(),
		// gin.Logger(),
		// gin.Recovery(),
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.ForceUA(),
	)
}

func setup404Handler(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			c.String(http.StatusNotFound, "404")
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "not found.",
			})
		}
	})
}
