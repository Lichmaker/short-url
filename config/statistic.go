package config

import "shorturl/pkg/config"

func init() {
	config.Add("statistic", func() map[string]interface{} {
		return map[string]interface{}{
			"enable": config.Env("STATISTIC_ENABLE", 0), // 是否启动统计。为0时不会触发统计写入
		}
	})
}
