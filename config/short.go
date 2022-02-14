// Package config 站点配置信息
package config

import "shorturl/pkg/config"

func init() {
	config.Add("short", func() map[string]interface{} {
		return map[string]interface{}{
			// 最大短链接缓存数量
			"cache_max": config.Env("SHORT_CACHE_MAX", 1000),
		}
	})
}
