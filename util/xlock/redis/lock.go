package redis

import (
	"context"
	"github.com/Remember9/frame/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

type Lock struct {
	Rdb    redis.Cmdable
	Name   string
	Expire time.Duration
	isLock bool
}

func Newlock(rdb redis.Cmdable, name string, expire time.Duration) *Lock {
	return &Lock{
		Rdb:    rdb,
		Name:   "lock_" + name,
		Expire: expire,
	}
}

// Lock 非阻塞锁
func (r *Lock) Lock(ctx context.Context) bool {
	re, err := r.Rdb.SetNX(ctx, r.Name, time.Now().Unix(), r.Expire).Result()
	if err != nil {
		xlog.Error("lock err", xlog.String("xlock-lockName", r.Name), xlog.Any("xlock-Duration", r.Expire), xlog.FieldErr(err))
		return false
	}
	if re {
		r.isLock = true
	}
	return re
}

// BlockLock 带阻塞带锁
func (r *Lock) BlockLock(ctx context.Context, retryNum int, sleepTime time.Duration) bool {
	limit := 1
	if retryNum > 1 {
		limit = retryNum
	}
	for retry := 0; retry < limit; retry++ {
		if r.Lock(ctx) {
			return true
		}
		if retry < limit {
			time.Sleep(sleepTime)
		}
	}
	return false
}

// UnLock 解锁
func (r *Lock) UnLock(ctx context.Context) bool {
	_, err := r.Rdb.Del(ctx, r.Name).Result()
	if err != nil {
		xlog.Error("unlock err", xlog.String("xlock-lockName", r.Name), xlog.FieldErr(err))
		return false
	} else {
		r.isLock = false
		return true
	}
}

// ReNew 续期
func (r *Lock) ReNew(ctx context.Context) bool {
	if r.isLock {
		re, err := r.Rdb.Expire(ctx, r.Name, r.Expire).Result()
		if err != nil {
			xlog.Error("Expire lock err", xlog.String("xlock-lockName", r.Name), xlog.FieldErr(err))
			return false
		}
		return re
	}
	return false
}
