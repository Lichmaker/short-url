package v1

import (
	"net/http"
	"shorturl/pkg/response"
	shortCore "shorturl/pkg/short"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseAPIController
}

func (controller *IndexController) Go(c *gin.Context) {
	shortStr := c.Param("short")
	long, ok := shortCore.Get(shortStr)
	if !ok {
		response.Abort404(c)
	} else {
		c.Redirect(http.StatusFound, long)
	}
}
