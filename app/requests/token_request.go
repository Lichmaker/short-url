package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type TokenRequest struct {
	AppID string `json:"app_id" valid:"appid"`
}

func GetToken(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"appid": []string{"required", "len:32"},
	}
	messages := govalidator.MapData{
		"appid": []string{
			"required:appid必传",
			"len:appid长度错误",
		},
	}
	return validate(data, rules, messages)
}
