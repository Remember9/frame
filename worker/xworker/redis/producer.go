package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Producer struct {
	*redisMq
	MaxLen int64 //队列长度
	Approx bool  //近似长度
}

// NewProducer 新建生产者
//rdb redis客户端
//name 队列名称
//num 控制队列长度，0标示不控制
func NewProducer(rdb redis.Cmdable, name string, num int64) *Producer {
	sName := StreamName(name)
	approx := false
	if num > 0 {
		approx = true
	}
	return &Producer{
		redisMq: &redisMq{stream: sName, Rdb: rdb, originName: name},
		MaxLen:  num,
		Approx:  approx,
	}
}
func (r *Producer) Send(ctx context.Context, data map[string]interface{}) (rid string, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in Consume err:%#v", err)
			rid = ""
		}
	}()
	arg := &redis.XAddArgs{
		Stream: r.stream,
		MaxLen: r.MaxLen,
		Approx: r.Approx,
		Values: data,
	}
	return r.Rdb.XAdd(ctx, arg).Result()
}
func (r *Producer) SendSlice(ctx context.Context, data []string) (rid string, topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in Consume err:%#v", err)
			rid = ""
		}
	}()
	arg := &redis.XAddArgs{
		Stream: r.stream,
		MaxLen: r.MaxLen,
		Approx: r.Approx,
		Values: data,
	}
	return r.Rdb.XAdd(ctx, arg).Result()
}
