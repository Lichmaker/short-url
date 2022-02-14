package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type ShortRequest struct {
	Url string `json:"url" valid:"url"`
}

func GetShort(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"url": []string{"required", "url"},
	}
	messages := govalidator.MapData{
		"url": []string{
			"required:url必传",
			"url:传入不是一个正确的url",
		},
	}
	return validate(data, rules, messages)
}
