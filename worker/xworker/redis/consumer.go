package redis

import (
	"context"
	"encoding/json"
	"fmt"
	rLock "github.com/Remember9/frame/util/xlock/redis"
	"github.com/Remember9/frame/xlog"
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"time"
)

type Consumer struct {
	*redisMq
	groupName string
	consumer  string               // 消费者
	Func      func(b []byte) error // 接收json的处理方式
	Pool      *ants.Pool
}

// NewConsume 新建消费者
// rdb redis客户端
// name 队列名称
// pool 协程池
// f 处理消息方法,参数为队列消息，格式为字节数组
func NewConsume(rdb redis.Cmdable, name string, pool *ants.Pool, f func([]byte) error) *Consumer {
	sName := StreamName(name)
	c := &Consumer{
		redisMq:   &redisMq{stream: sName, Rdb: rdb, originName: name},
		groupName: sName,
		consumer:  consumerName(name),
		Pool:      pool,
		Func:      f,
	}
	return c
}

// Consume 处理池外部定义
func (r *Consumer) Consume(ctx context.Context) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in Consume err:%#v", err)
		}
	}()

	r.Rdb.XGroupCreateMkStream(ctx, r.stream, r.groupName, "0")

	for {
		readArg := &redis.XReadGroupArgs{
			Group:    r.groupName,
			Consumer: r.consumer,
			Streams:  []string{r.stream, ">"},
			// count is number of entries we want to read from redis
			Count: 1,
			// we use the block command to make sure if no entry is found we wait
			// until an entry is found
			// 0表示一直阻塞到消息来
			// 1小时超时
			Block: time.Hour,
		}
		data, err := r.Rdb.XReadGroup(ctx, readArg).Result()
		// xlog.Infof("get mq data----Group:%s,Consumer:%s,Streams:%s,data:%#v", r.GroupName, r.Consumer, r.Stream, data)
		// 错误休眠2秒后重试

		if err != nil {
			xlog.Info("XReadGroup 获取信息异常,2秒后重试", xlog.FieldErr(err), xlog.Any("redis-mq-arg", readArg))
			time.Sleep(time.Second * 2)
			continue
		}
		if data == nil || len(data) == 0 || len(data[0].Messages) == 0 {
			xlog.Info("获取data容量为空,retry now")
			continue
		}
		// 投入协程池处理
		if err = r.addPoll(ctx, data[0]); err != nil {
			xlog.Error("任务协程池异常,2秒后重试", xlog.FieldErr(err))
			time.Sleep(time.Second * 2)
			continue
		}
	}
}
func (r *Consumer) addPoll(ctx context.Context, da redis.XStream) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in redis Consumer addPoll err:%#v", err)
		}
	}()
	err := r.Pool.Submit(func() {
		st := time.Now()
		str, err := json.Marshal(da.Messages[0].Values)
		if err != nil {
			xlog.Error("json.Marshal data Err", xlog.Any("redis-mq-data", da), xlog.String("redis-mq-stream", r.stream), xlog.String("redis-mq-groupName", r.groupName), xlog.FieldErr(err))
		}
		if err := r.Func(str); err != nil {
			// 错误处理
			cost := time.Since(st).Milliseconds()
			xlog.Error("consume data Err", xlog.Any("redis-mq-data", da), xlog.String("redis-mq-stream", r.stream), xlog.String("redis-mq-groupName", r.groupName), xlog.Any("redis-mq-cost", cost), xlog.FieldErr(err))
		} else {
			// 正常时ack
			cost := time.Since(st).Milliseconds()
			ackRes, ackErr := r.Rdb.XAck(ctx, da.Stream, r.groupName, da.Messages[0].ID).Result()
			xlog.Info("consume data success", xlog.Any("redis-mq-data", da), xlog.String("redis-mq-stream", r.stream), xlog.String("redis-mq-groupName", r.groupName), xlog.Any("redis-mq-cost", cost),
				xlog.Any("redis-mq-ackRes", fmt.Sprintf("%v", ackRes)), xlog.Any("redis-mq-ackErr", fmt.Sprintf("%v", ackErr)))
		}
	})

	if err != nil {
		return err
	}
	return nil
}

// TickerPending 定时任务执行pending
// expTime 过期时间(当时间小于等于0时设置为默认 过期时间为1天)
// perTime 执行周期(当时间小于等于0时设置为默认 每15分钟)
// startRun 是否开始就先执行一次
func (r *Consumer) TickerPending(ctx context.Context, expTime time.Duration, perTime time.Duration, startRun bool, MaxLen int64) {
	if expTime <= 0 {
		expTime = 24 * time.Hour
	}
	if perTime <= 0 {
		perTime = 15 * time.Minute
	}
	if startRun {
		r.DispatchPending(expTime, MaxLen)
	}
	ticker := time.NewTicker(perTime)
	go func() {
		for {
			select {
			case <-ticker.C:
				r.DispatchPending(expTime, MaxLen)
			case <-ctx.Done():
				ticker.Stop()
				xlog.Info("关闭TickerPending")
				return
			}
		}
	}()
}

// DispatchPending 对超时pending进行回收
// 每一个pending获取一次消息（避免直接range一批因范围问题导致消息丢失）
func (r *Consumer) DispatchPending(expTime time.Duration, MaxLen int64) {
	ctx := context.Background()
	defer func() {
		if err := recover(); err != nil {
			xlog.Error("dispatchPending panic", xlog.Any("redis-mqPending-err", err), xlog.String("redis-mqPending-stream", r.stream),
				xlog.String("redis-mqPending-groupName", r.groupName))
		}
	}()
	r.Rdb.XGroupCreateMkStream(ctx, r.stream, r.groupName, "0")
	// 加锁，未拿到锁就跳过
	lock := rLock.Newlock(r.Rdb, "Pending_"+r.stream, time.Second*300)
	if !lock.Lock(ctx) {
		return
	}
	defer lock.UnLock(ctx)
	xlog.Info("dispatchPending start", xlog.String("redis-mqPending-stream", r.stream), xlog.String("redis-mqPending-groupName", r.groupName),
		xlog.Int64("redis-mqPending-setLength", MaxLen))
	defer xlog.Info("dispatchPending end", xlog.String("redis-mqPending-stream", r.stream), xlog.String("redis-mqPending-groupName", r.groupName))
	pInof, err := r.Rdb.XPending(ctx, r.stream, r.groupName).Result()
	if err != nil {
		xlog.Error("dispatchPending error XPending Err", xlog.FieldErr(err), xlog.String("redis-mqPending-stream", r.stream),
			xlog.String("redis-mqPending-groupName", r.groupName))
		return
	}
	if MaxLen < 0 {
		MaxLen = 0
	}
	producer := NewProducer(r.Rdb, r.originName, MaxLen)
	if pInof.Count > 0 {
		for cu, num := range pInof.Consumers {
			xExt := &redis.XPendingExtArgs{
				Stream:   r.stream,
				Group:    r.groupName,
				Idle:     expTime, // 15分钟未消费的
				Start:    "-",
				End:      "+",
				Count:    num,
				Consumer: cu,
			}
			list, err := r.Rdb.XPendingExt(ctx, xExt).Result()
			if err != nil {
				xlog.Error("dispatchPending XPendingExt Err", xlog.FieldErr(err), xlog.Any("redis-mqPending-PendingArg", xExt),
					xlog.String("redis-mqPending-stream", r.stream), xlog.String("redis-mqPending-groupName", r.groupName))
				continue
			}
			if len(list) == 0 {
				continue
			}
			ids := make([]string, 0)
			pendingIds := make(map[string]struct{})
			for _, v := range list {
				ids = append(ids, v.ID)
				pendingIds[v.ID] = struct{}{}
			}
			for _, pendid := range ids {
				rangeList, err2 := r.Rdb.XRangeN(ctx, r.stream, pendid, pendid, 1).Result()
				if err2 != nil {
					xlog.Error("dispatchPending XRangeN Err", xlog.FieldErr(err2), xlog.String("redis-mqPending-stream", r.stream),
						xlog.String("redis-mqPending-start", pendid), xlog.String("redis-mqPending-end", pendid),
						xlog.Any("redis-mqPending-count", 1))

					continue
				}
				if len(rangeList) == 0 {
					// 无法从range中取值，直接ack
					newId := "0"
					raValue := map[string]interface{}{"last_id": pendid}
					ack, ackerr := r.Rdb.XAck(ctx, r.stream, r.groupName, pendid).Result()
					xlog.Info("dispatchPending", xlog.String("redis-mqPending-stream", r.stream),
						xlog.String("redis-mqPending-id", newId), xlog.Any("redis-mqPending-message", raValue),
						xlog.Any("redis-mqPending-ack", fmt.Sprintf("%v", ack)),
						xlog.Any("redis-mqPending-ackErr", fmt.Sprintf("%v", ackerr)))
				} else {
					// 可以从range取值
					ra := rangeList[0]
					rav := ra.Values
					rav["last_id"] = ra.ID
					r.Rdb.Pipeline()
					newId, err := producer.Send(ctx, rav)
					if err != nil {
						xlog.Error("dispatchPending send Err", xlog.FieldErr(err), xlog.String("redis-mqPending-stream", r.stream),
							xlog.Any("redis-mqPending-data", rav))
						continue
					}
					ack, ackerr := r.Rdb.XAck(ctx, r.stream, r.groupName, ra.ID).Result()
					xlog.Info("dispatchPending", xlog.String("redis-mqPending-stream", r.stream),
						xlog.String("redis-mqPending-id", newId), xlog.Any("redis-mqPending-message", ra.Values),
						xlog.Any("redis-mqPending-ack", fmt.Sprintf("%v", ack)),
						xlog.Any("redis-mqPending-ackErr", fmt.Sprintf("%v", ackerr)))
				}

			}
		}
	}
	return
}

/*
*
获取消费者名称（根据需要加处理）
*/
func consumerName(name string) string {
	return name
}
