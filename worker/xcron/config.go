package xcron

import (
	"esfgit.leju.com/golang/frame/client/redis"
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/xlog"
	"fmt"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	"runtime"
	"time"
)

func DefaultConfig() Config {
	return Config{
		wrappers:        []JobWrapper{},
		WithSeconds:     false, // 是否使用秒级cron
		ImmediatelyRun:  false, // 是否立即运行
		ConcurrentDelay: -1,    // 任务并发时是否延迟运行
		DistributedTask: true,  // 是否并发加锁
	}
}

// Config ...
type Config struct {
	WithSeconds     bool
	ConcurrentDelay int
	ImmediatelyRun  bool

	wrappers []JobWrapper
	xparser  cron.Parser

	// Distributed task
	DistributedTask bool
	rdb             *redis.Redis
}

// Build ...
func Build(name string) *Cron {
	var cronConfig = DefaultConfig()
	if err := config.UnmarshalKey(name, &cronConfig); err != nil {
		xlog.Panic("cron parse config panic", logMod, xlog.FieldErr(err), xlog.String("key", name), xlog.Any("value", cronConfig))
	}

	if cronConfig.WithSeconds {
		cronConfig.xparser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	} else { // default parser
		cronConfig.xparser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	if cronConfig.ConcurrentDelay > 0 { // 延迟
		cronConfig.wrappers = append(cronConfig.wrappers, delayIfStillRunning())
	} else if cronConfig.ConcurrentDelay < 0 { // 跳过
		cronConfig.wrappers = append(cronConfig.wrappers, skipIfStillRunning())
	}

	return newCron(&cronConfig)
}

type wrappedJob struct {
	NamedJob

	distributedTask bool
	rdb             *redis.Redis
	lockExpire      time.Duration // 需要大于执行时间
}

// Run ...
func (wj wrappedJob) Run() {
	if wj.distributedTask {
		ctx := context.Background()
		lockKey := "xcron:lock:" + wj.Name()
		rx := wj.rdb.SetNx(ctx, lockKey, 1, wj.lockExpire)
		if !rx {
			xlog.Info("locked", logMod, xlog.String("name", wj.Name()))
			return
		}
		defer wj.rdb.Del(ctx, lockKey)
	}
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	var fields = []xlog.Field{logMod, xlog.String("name", wj.Name())}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, xlog.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, xlog.String("err", err.Error()), xlog.Duration("cost", time.Since(beg)))
			xlog.Error("run", fields...)
		} else {
			xlog.Info("run", fields...)
		}
	}()

	return wj.NamedJob.Run()
}
