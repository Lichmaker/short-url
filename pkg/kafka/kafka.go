package kafka

import (
	"shorturl/pkg/config"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func GetClientConfig() *sarama.Config {
	clientConfig := sarama.NewConfig()
	clientConfig.ClientID = config.GetString("kafka.clientid") // client 名称，用于给broker中日志记录
	clientConfig.Version = sarama.V3_1_0_0                     // kafka server的版本号
	clientConfig.Producer.Return.Successes = true              // sync必须设置这个
	clientConfig.Producer.RequiredAcks = sarama.WaitForAll     // 也就是等待foolower同步，才会返回
	clientConfig.Producer.Return.Errors = true
	clientConfig.Consumer.Return.Errors = true
	clientConfig.Metadata.Full = false                                                // 不用拉取全部的信息
	clientConfig.Consumer.Offsets.AutoCommit.Enable = true                            // 自动提交偏移量，默认开启，说时候，我没找到手动提交。
	clientConfig.Consumer.Offsets.AutoCommit.Interval = time.Second                   // 这个看业务需求，commit提交频率，不然容易down机后造成重复消费。
	clientConfig.Consumer.Offsets.Initial = sarama.OffsetOldest                       // 从最开始的地方消费，业务中看有没有需求，新业务重跑topic。
	clientConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin // rb策略，默认就是range
	return clientConfig
}

func GetAddress() []string {
	addrString := config.GetString("kafka.address")
	return strings.Split(addrString, ",")
}

func GetProducer() (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(GetAddress(), GetClientConfig())
	if err != nil {
		return nil, err
	}
	return producer, nil
}
