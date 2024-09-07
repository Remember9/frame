package kafka

import "github.com/Shopify/sarama"

type Base struct {
	topic      string   //topic名
	originName string   //未加前缀队列名
	addrs      []string //["172.0.0.1:9092"]
	Config     *sarama.Config
}
