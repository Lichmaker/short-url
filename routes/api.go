package routes

import (
	controllers "shorturl/app/http/controllers/api/v1"
	"shorturl/app/http/controllers/api/v1/auth"
	"shorturl/app/http/middlewares"
	"shorturl/pkg/config"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}
	{
		dc := new(controllers.DemoController)
		v1.GET("hello", dc.HelloWorld)

		authGroup := v1.Group("/auth")
		{
			tokenController := new(auth.TokenController)
			authGroup.POST("/get-token", middlewares.GuestJWT(), tokenController.Get) // todo 现在没有做appid，可以随便生成，不加limit
		}

		urlController := new(controllers.UrlController)
		v1.POST("/short", middlewares.AuthJWT(), middlewares.LimitPerRoute("60-M"), urlController.Short)
	}

	indexController := new(controllers.IndexController)
	r.GET("/:short", indexController.Go)
}
