package redis

import (
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/xlog"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func init() {
	err := config.InitTest()
	if err != nil {
		panic(err)
	}
	err = xlog.Build()
	if err != nil {
		panic(err)
	}
}

func TestRedis(t *testing.T) {
	redisClient := Build("redis")
	st := redisClient.Stub().PoolStats()
	t.Logf("running status %+v", st)
	err := redisClient.Close()
	if err != nil {
		t.Errorf("redis close failed:%v", err)
	}
	st = redisClient.Stub().PoolStats()
	t.Logf("close status %+v", st)
}

func TestCmds(t *testing.T) {
	ctx := context.Background()

	r := Build("redis")
	if err := r.Stub().Ping(ctx).Err(); err != nil {
		xlog.Error("start redis", xlog.Any("err", err))
	}

	setRes := r.Set(ctx, "demo_key", "demo_value", time.Second*50)
	t.Logf("redis set string %v", setRes)

	getRes := r.Get(ctx, "demo")
	t.Logf("redis get string %v", getRes)
}
