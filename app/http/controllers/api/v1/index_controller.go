package v1

import (
	"net/http"
	"shorturl/pkg/response"
	shortCore "shorturl/pkg/short"
	"shorturl/pkg/statistic"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseAPIController
}

func (controller *IndexController) Go(c *gin.Context) {
	shortStr := c.Param("short")
	model, ok := shortCore.Get(shortStr)
	if !ok {
		response.Abort404(c)
	} else {
		// 计数统计
		statistic.Enqueue(model)
		c.Redirect(http.StatusFound, model.Long)
	}
}
