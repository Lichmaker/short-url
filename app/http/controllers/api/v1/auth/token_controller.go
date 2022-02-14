package auth

import (
	v1 "shorturl/app/http/controllers/api/v1"
	"shorturl/app/requests"
	"shorturl/pkg/jwt"
	"shorturl/pkg/logger"
	"shorturl/pkg/response"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	v1.BaseAPIController
}

func (controller *TokenController) Get(c *gin.Context) {
	// 接受参数appkey，只有对应的appkey才可以生成token
	request := &requests.TokenRequest{}
	if !requests.Validate(c, request, requests.GetToken) {
		return
	}

	// var abPath string
	// _, filename, _, ok := runtime.Caller(0)
	// if ok {
	// 	abPath = path.Dir(filename)
	// }
	// ospath, _ := os.Getwd()

	logger.DebugString("get_token", "传入的appid", request.AppID)
	// todo 验证appid是否存在

	token := jwt.NewJWT().IssueToken(request.AppID)
	response.Data(c, gin.H{
		"token": token,
		// "p":      app.AbsolutePath,
		// "abpath": abPath,
		// "work":   ospath,
	})
}
