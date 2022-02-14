// Package app 应用信息
package app

import (
	"shorturl/pkg/config"
	"strings"
	"time"
)

// 会记录项目可执行文件的绝对路径
var AbsolutePath string

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}

// URL 传参 path 拼接站点的 URL
func URL(path string) string {
	return strings.TrimSuffix(config.GetString("app.url"), "/") + "/" + strings.TrimPrefix(path, "/")
}

// V1URL 拼接带 v1 标示 URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}
