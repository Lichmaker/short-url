package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"shorturl/pkg/app"
	cmdhelper "shorturl/pkg/cmd-helper"
	"shorturl/pkg/config"
	"shorturl/pkg/console"
	"shorturl/pkg/kafka"
	"shorturl/pkg/logger"
	"shorturl/pkg/statistic"
	"shorturl/pkg/traceid"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var CmdConsumer = &cobra.Command{
	Use:   "consumer",
	Short: "消费者",
	Run:   runConsumer,
}

func init() {
	// 注册子命令
	CmdConsumer.AddCommand(
		CmdShutdown,
	)
}

// 传入启动数量
var RunConsumerCount string

// 标记为启动实例
var RunInstance bool

// 消费者组client
var client sarama.ConsumerGroup

// 消费者组client的上下文，用于暂停、结束
var ctx context.Context

// waitgroup 用于等待信号
var wg *sync.WaitGroup

// 结束消费
var cancel context.CancelFunc

func runConsumer(cmd *cobra.Command, args []string) {
	console.Info("starting...")

	if RunInstance {
		// 不存在，启动
		i := args[0]
		handler(cast.ToString(i))
	} else {
		count := cast.ToInt(RunConsumerCount)
		if count == 0 {
			console.Error("请指定启动数量")
			return
		}

		createCount := 0
		for i := 1; i <= count; i++ {
			// 先查询是否已存在
			processName := fmt.Sprintf("%s consumer -i %d", config.GetString("app.process_name"), i)
			psCmd := `ps aux | awk '/` + processName + `/ && !/awk/ && !/shutdown/ {print $1}'`
			pid, err := cmdhelper.RunCommand(psCmd)
			logger.LogIf(err)
			if len(pid) <= 0 {
				// 不存在，启动
				createCount++
				handler(cast.ToString(i))
			} else {
				// 已存在，跳过
				continue
			}
		}
		console.Success(fmt.Sprintf("成功启动消费者 %d 个", createCount))
	}

}

func handler(instanceNumber string) {
	traceid.Boot("")
	pidfilename := fmt.Sprintf("shorturlkafka-%s.pid", instanceNumber)
	logfilename := fmt.Sprintf("shorturlkafka-%s.log", instanceNumber)
	// daemon
	cntxt := &daemon.Context{
		PidFileName: pidfilename,
		PidFilePerm: 0644,
		LogFileName: logfilename,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{app.AbsolutePath + "/" + config.GetString("app.process_name"), "consumer", "-i", instanceNumber},
	}

	d, err := cntxt.Reborn()
	logger.LogIf(err)
	if d != nil {
		return
	}
	defer cntxt.Release()

	group := statistic.KAFKA_COUSUMER_GROUP
	topic := statistic.KAFKA_TOPIC
	addr := kafka.GetAddress()
	logger.InfoJSON("consumer", "address", addr)
	client, err = sarama.NewConsumerGroup(addr, group, kafka.GetClientConfig()) // broker_ip，消费者组(broker记录偏移量)，kafka 配置设置
	if err != nil {
		logger.LogIf(err)
		panic(err)
	}

	wg = &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel = context.WithCancel(context.Background())

	go worker(client, topic)

	// 注册信号
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			console.Warning("consumer term: 接收到上下文")

			keepRunning = false
		case <-sigterm:
			console.Warning("consumer term: 接收到信号")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		logger.LogIf(err)
		panic(err)
	}

	console.Exit("消费者实例结束")
}

func worker(client sarama.ConsumerGroup, topic string) {
	defer wg.Done()
WORKERLOOP:
	for { // for循环的目的是因为存在重平衡，他会重新启动

		if ctx.Err() != nil {
			break WORKERLOOP
		}
		console.Info("consume 开始...")
		handler := new(ConsumerGroupHandler)                 // 必须传递一个handler
		err := client.Consume(ctx, []string{topic}, handler) // consume 操作，死循环。exampleConsumerGroupHandler的ConsumeClaim不允许退出，也就是操作到完毕。
		if err != nil {
			logger.LogIf(err)
			panic(err)
		}
	}
}

type ConsumerGroupHandler struct{}

var _ sarama.ConsumerGroupHandler = ConsumerGroupHandler{}

func (h ConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	console.Info("kafka连接完成")
	return nil
}

func (h ConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	console.Info("kafka消费 cleanup")
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
FORLOOP:
	for msg := range claim.Messages() { // 接受topic消息
		if ctx.Err() != nil {
			break FORLOOP
		}

		// fmt.Printf("[Consumer] Message topic:%q partition:%d offset:%d add:%d\n", msg.Topic, msg.Partition, msg.Offset, claim.HighWaterMarkOffset()-msg.Offset)
		logger.DebugJSON("consumer", "ConsumeClaim", gin.H{
			"topic":     msg.Topic,
			"partition": msg.Partition,
			"offset":    msg.Offset,
			"value":     msg.Value,
		})
		err := statistic.Handle(string(msg.Value))
		if err != nil {
			logger.LogIf(err)
		} else {
			sess.MarkMessage(msg, "") // 必须设置这个，不然你的偏移量无法提交。
		}
	}
	return nil
}
