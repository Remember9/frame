package kafka

import (
	"context"
	"esfgit.leju.com/golang/frame/xlog"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type Consumer struct {
	*Base
	groupName    string
	Func         func(b []byte) error //接收json的处理方式
	Client       sarama.ConsumerGroup
	groupHandler *GroupHandler
}

// NewConsumer 创建消费者
//创建消费者后可以调整config相关参数
//addrs地址数组["172.0.0.1:9092","172.0.0.1:9093"]
//version kafka版本号
//f 处理方法
func NewConsumer(addrs []string, topic string, version string, f func(b []byte) error) *Consumer {
	p := &Consumer{}
	p.groupName = topic
	p.Func = f
	p.groupHandler = &GroupHandler{
		Func: f,
	}
	config := sarama.NewConfig()
	kversion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		panic(err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config.Version = kversion
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	p.Base = &Base{
		topic:      topic,
		originName: topic,
		addrs:      addrs,
		Config:     config,
	}
	return p
}

func (p *Consumer) Init() (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka Consumer init err:%#v", err)
		}
	}()
	client, err := sarama.NewConsumerGroup(p.addrs, p.groupName, p.Config)
	if err != nil {
		return err
	}
	p.Client = client
	return nil
}

func (p *Consumer) Consume(ctx context.Context) (topErr error) {
	defer func() {
		if err := p.Client.Close(); err != nil {
			xlog.Info("Error closing client", xlog.FieldErr(err), xlog.String("topic", p.topic), xlog.Any("addr", p.addrs))
		}
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka Consumer init err:%#v", err)
		}
	}()

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		//此处可同时消费多个topic队列，但为了后续处理逻辑简单只指定消费一个队列
		if err := p.Client.Consume(ctx, []string{p.topic}, p.groupHandler); err != nil {
			xlog.Info("获取消费信息异常，2秒后重试", xlog.FieldErr(err), xlog.String("topic", p.topic), xlog.Any("addr", p.addrs))
			time.Sleep(time.Second * 2)
		}
		// check if context was cancelled, signaling that the consumer should stop
		//如果外部使用ctx的cancal，则退出消费
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// GroupHandler  represents a Sarama consumer group consumer
//在分区变更时会触发GroupHandler退出重启
type GroupHandler struct {
	Func func(b []byte) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *GroupHandler) Setup(ss sarama.ConsumerGroupSession) error {
	xlog.Info("kafka 消费者启动", xlog.Any("Claims", ss.Claims()), xlog.String("MemberID", ss.MemberID()),
		xlog.Any("GenerationID", ss.GenerationID()))
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *GroupHandler) Cleanup(ss sarama.ConsumerGroupSession) error {
	xlog.Info("kafka 消费者退出", xlog.Any("Claims", ss.Claims()), xlog.String("MemberID", ss.MemberID()),
		xlog.Any("GenerationID", ss.GenerationID()))
	return nil
}

func (c *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (topErr error) {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	//这块原package对于一个topic每一个分区消费已单独起了一个协程，consumer_group.go的659行
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka ConsumeClaim err:%#v", err)
		}
	}()
	for {
		select {
		case message := <-claim.Messages():
			func(message *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
				defer session.MarkMessage(message, "")
				err := c.Func(message.Value)
				if err != nil {
					xlog.Info("consume data error", xlog.FieldErr(err), xlog.String("data", string(message.Value)), xlog.Any("message", message))
				}
			}(message, session)
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			//session通知的取消事件要退出
			//需要保留此监听
			return nil
		}
	}
}

/*// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *GroupHandler) ConsumeClaim_whthPool(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (topErr error) {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	//这块原package对于一个topic每一个分区消费已单独起了一个协程，consumer_group.go的659行
	//！！！！！kafka使用协程池对于commit每次的offset会默认之前的offict都成功，似乎不能使用协程池
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka ConsumeClaim err:%#v", err)
		}
	}()
	for {
		select {
		case message := <-claim.Messages():
			err := c.addPool(message, session)
			if err != nil {
				xlog.Info("任务协程池异常", xlog.FieldErr(err), xlog.Any("message", message))
			}
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			//session通知的取消事件要退出
			//需要保留此监听
			return nil
		}
	}
}
func (c *GroupHandler) addPool(message *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka GroupHandler addPoll err:%#v", err)
		}
	}()
	err := c.Pool.Submit(func() {
		err := c.Func(message.Value)
		if err != nil {
			xlog.Info("consume data error", xlog.FieldErr(err), xlog.String("data", string(message.Value)), xlog.Any("message", message))
		} else {
			session.MarkMessage(message, "")
			//session.Commit() //默认每一秒自动提交带消费标记信息
			xlog.Info("consume data success", xlog.String("data", string(message.Value)), xlog.Any("message", message))
		}
	})
	return err
}*/
