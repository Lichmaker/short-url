package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"shorturl/bootstrap"
	btsConfig "shorturl/config"
	"shorturl/pkg/config"
	"shorturl/pkg/helpers"
	"shorturl/pkg/traceid"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 基础服务boot
	absolutePath, _ := os.Getwd()
	absolutePathSplit := strings.Split(absolutePath, "/")
	finalPath := strings.Join(absolutePathSplit[:len(absolutePathSplit)-6], "/")
	envPath := finalPath + "/.env"
	traceid.Boot("")
	config.InitConfig("", envPath)
	// 初始化所有的config
	btsConfig.Initialize()
	bootstrap.SetupLogger()
	bootstrap.SetupDB()
	bootstrap.SetupRedis()
	bootstrap.SetupCache()
}

func TestGet(t *testing.T) {
	// 定义测试用例
	testsCases := []struct {
		appid string
	}{
		{
			appid: helpers.RandomString(32),
		},
		{
			appid: helpers.RandomString(32),
		},
	}

	// http服务的bootstrap
	router := gin.New()
	// 初始化路由绑定
	bootstrap.SetupRoute(router)

	for _, caseItem := range testsCases {
		t.Run("生成JWT", func(t *testing.T) {
			requestJson, e := json.Marshal(struct {
				App_id string `json:"app_id"`
			}{caseItem.appid})
			fmt.Println(string(requestJson))
			assert.Nil(t, e)
			// 建立一个HTTP请求
			req := httptest.NewRequest(
				"POST",                                 // 请求方法
				"/api/v1/auth/get-token",               // 请求URL
				strings.NewReader(string(requestJson)), // 请求参数
			)
			req.Header.Add("User-Agent", "golang-test")
			req.Header.Add("Content-Type", "application/json")

			// mock一个响应记录器
			w := httptest.NewRecorder()

			// 让server端处理mock请求并记录返回的响应内容
			router.ServeHTTP(w, req)

			// 校验状态码是否符合预期
			assert.Equal(t, http.StatusOK, w.Code)

			// 解析并检验响应内容是否复合预期
			var resp struct {
				Data struct {
					Token string `json:"token"`
				} `json:"data"`
				Success bool `json:"success"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			fmt.Println(resp)
			assert.Nil(t, err)
			assert.Equal(t, 261, len(resp.Data.Token))
		})
	}
}
