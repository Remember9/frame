package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Remember9/frame/util/xpool"
	"github.com/Remember9/frame/xlog"
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"time"
)

type Pending struct {
	ExpTime  time.Duration // 过期时间
	PerTime  time.Duration // 回收触发时间间隔
	StartRun bool          // 是否在开始就执行一次回收操作
	MaxLen   int64         // 队列长度
}
type Worker struct {
	Name    string               // 队列名称
	Handle  func(m []byte) error // 处理方法
	Pending *Pending             // 如果为nil则不启动pending回收
	Rdb     redis.Cmdable        // redis client
	PoolNum uint                 // 线程池数，如果不传则走全局池，数目从配置文件读取
	pool    *ants.Pool
}

func (r *Worker) Init() {
	var err error
	if r.PoolNum == 0 {
		r.pool, _, err = xpool.GetAntsPool("workerGlobal")
	} else {
		r.pool, err = xpool.NewAntsPool(r.Name, int(r.PoolNum), xpool.WithAntsLogger())
	}
	if err != nil {
		xlog.Panic(err.Error(), xlog.FieldErr(err))
		// panic(err)
	}
}

func (r *Worker) GetName() string {
	return r.Name
}
func (r *Worker) Serve(ctx context.Context) (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("%s:panic in redis worker serve err:%#v", r.Name, err)
		}
	}()
	if r.pool == nil {
		return errors.New("WorkerPool 没有初始化")
	}
	c := NewConsume(r.Rdb, r.Name, r.pool, r.Handle)
	if r.Pending != nil {
		// ------定时回收过期未处理的pending消息-----start
		c.TickerPending(ctx, r.Pending.ExpTime, r.Pending.PerTime, r.Pending.StartRun, r.Pending.MaxLen)
		// ------定时回收过期未处理的pending消息-----end
	}
	// 无需修改
	if err := c.Consume(ctx); err != nil {
		return err
	}
	return nil
}
