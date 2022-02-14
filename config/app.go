package config

import "shorturl/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			"name":     config.Env("APP_NAME", "myApp"),
			"env":      config.Env("APP_ENV", "production"),
			"debug":    config.Env("APP_DEBUG", false),
			"port":     config.Env("APP_PORT", "8080"),
			"key":      config.Env("APP_KEY", "abc"),
			"url":      config.Env("APP_URL", "http://127.0.0.1:8080"),
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
			// API 域名，未设置的话所有 API URL 加 api 前缀，如 http://domain.com/api/v1/users
			"api_domain":   config.Env("API_DOMAIN"),
			"process_name": config.Env("APP_PROCESS_NAME", "shorturl"),
		}
	})
}
