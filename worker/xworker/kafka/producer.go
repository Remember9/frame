package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
)

type Producer struct {
	*Base
	syncProducer sarama.SyncProducer
}

// NewProducer
//创建生产者后需要调用init初始化，初始化之前可以调整config相关参数
func NewProducer(addrs []string, topic string) *Producer {
	p := &Producer{}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	p.Base = &Base{
		topic:      topic,
		originName: topic,
		addrs:      addrs,
		Config:     config,
	}
	return p
}

// SyncInit 初始化同步模式生产者
func (p *Producer) SyncInit() (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka NewSyncProducer err:%#v", err)
		}
	}()
	producer, err := sarama.NewSyncProducer(p.addrs, p.Config)
	if err != nil {
		return err
	}
	p.syncProducer = producer
	return nil
}
func (p *Producer) SyncSend(data []byte) (partition int32, offset int64, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka Send err:%#v", err)
		}
	}()
	if p.syncProducer == nil {
		return 0, 0, errors.New("kafka producer not init")
	}
	return p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	})
}

// SyncClose 关闭方法在系统退出时可以调用
func (p *Producer) SyncClose() (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in kafka Close err:%#v", err)
		}
	}()
	return p.syncProducer.Close()
}
