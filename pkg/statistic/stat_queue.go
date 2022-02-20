package statistic

import (
	"shorturl/pkg/kafka"
	"shorturl/pkg/logger"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

func kafkaEnqueue(shortstr string) {
	producer, err := kafka.GetProducer()
	if err != nil {
		logger.LogIf(err)
		return
	}
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: KAFKA_TOPIC,
		Value: sarama.StringEncoder(shortstr),
	})
	logger.LogIf(err)
	logger.InfoJSON("statistic", "kafka", gin.H{
		"value":     shortstr,
		"partition": partition,
		"offset":    offset,
	})
}
