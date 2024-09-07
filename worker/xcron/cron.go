package xcron

import (
	"github.com/Remember9/frame/client/redis"
	"github.com/Remember9/frame/util/xstring"
	"github.com/Remember9/frame/xlog"
	"github.com/robfig/cron/v3"
	"sync/atomic"
	"time"
)

var (
	// Every ...
	Every = cron.Every
	// NewParser ...
	NewParser = cron.NewParser
	// NewChain ...
	NewChain = cron.NewChain
	// WithSeconds ...
	WithSeconds = cron.WithSeconds
	// WithParser ...
	WithParser = cron.WithParser
	// WithLocation ...
	WithLocation = cron.WithLocation
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Entry ...
	Entry = cron.Entry
	// Schedule ...
	Schedule = cron.Schedule
	// Parser ...
	Parser = cron.Parser
	// Option ...
	Option = cron.Option
	// Job ...
	Job = cron.Job
	// NamedJob ..
	NamedJob interface {
		Run() error
		Name() string
	}
)

type FuncJob func() error

func (f FuncJob) Run() error { return f() }

func (f FuncJob) Name() string {
	return xstring.FunctionName(f)
}

type Cron struct {
	*Config
	*cron.Cron
}

func newCron(config *Config) *Cron {
	c := &Cron{
		Config: config,
		Cron: cron.New(
			cron.WithParser(config.xparser),
			cron.WithChain(config.wrappers...),
			cron.WithLogger(&wrappedLogger{}),
		),
	}
	return c
}

// WithRdb distributedTask=true 必须使用
func (c *Cron) WithRdb(rdb *redis.Redis) *Cron {
	c.Config.rdb = rdb
	return c
}

func (c *Cron) Schedule(schedule Schedule, job NamedJob, lockExpire time.Duration) EntryID {
	if c.ImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}

	xlog.Info("add job", logMod, xlog.String("name", job.Name()))
	return c.Cron.Schedule(schedule, &wrappedJob{
		NamedJob:        job,
		distributedTask: c.DistributedTask,
		rdb:             c.rdb,
		lockExpire:      lockExpire,
	})
}

func (c *Cron) AddJob(spec string, cmd NamedJob, lockExpire time.Duration) (EntryID, error) {
	schedule, err := c.xparser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd, lockExpire), nil
}

func (c *Cron) AddFunc(spec string, cmd func() error, lockExpire time.Duration) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd), lockExpire)
}

func (c *Cron) Run() error {
	xlog.Info("run xcron", logMod, xlog.Int("number of scheduled jobs", len(c.Cron.Entries())))
	c.Cron.Run()
	return nil
}

func (c *Cron) Stop() error {
	_ = c.Cron.Stop()
	return nil
}

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}
