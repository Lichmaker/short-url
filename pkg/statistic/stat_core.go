package statistic

import (
	"shorturl/app/models/short"
	"shorturl/pkg/config"
	"shorturl/pkg/logger"
)

func Enable() bool {
	return config.GetInt("statistic.enable") == 1
}

func Enqueue(model short.Short) {
	if !Enable() {
		return
	}
	logger.DebugString("statistic", "Enqueue", model.Short)
	kafkaEnqueue(model.Short)
}
