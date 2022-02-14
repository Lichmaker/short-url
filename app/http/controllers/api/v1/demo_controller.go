package v1

import (
	"shorturl/pkg/response"

	"github.com/gin-gonic/gin"
)

type DemoController struct {
	BaseAPIController
}

func (controller *DemoController) HelloWorld(c *gin.Context) {
	response.Success(c)
}
