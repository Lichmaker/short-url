package v1

import (
	"shorturl/app/requests"
	"shorturl/pkg/response"
	shortCore "shorturl/pkg/short"

	"github.com/gin-gonic/gin"
)

type UrlController struct {
	BaseAPIController
}

func (controller *UrlController) Short(c *gin.Context) {
	request := &requests.ShortRequest{}
	if !requests.Validate(c, request, requests.GetShort) {
		return
	}

	//处理缩短
	shortUrl, ok := shortCore.Generate(request.Url)
	if !ok {
		response.Abort500(c)
	} else {
		response.Data(c, gin.H{
			"url":      shortUrl,
			"long_url": request.Url,
		})
	}
}
