package config

import "shorturl/pkg/config"

func init() {
	config.Add("kafka", func() map[string]interface{} {
		return map[string]interface{}{
			"clientid": config.Env("KAFKA_CLIENT_ID", "shorturl-kafka"),
			"address": config.Env("KAFKA_ADDRESS"),
		}
	})
}
