package xworker

import (
	"context"
	"errors"
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/util/xpool"
	"esfgit.leju.com/golang/frame/xlog"
	"fmt"
)

type ConsumerFunc interface {
	Serve(context.Context) error
	GetName() string
	Init()
}

type WorkerBasic struct {
	consumers  []ConsumerFunc
	CancelFunc context.CancelFunc
}

func (w *WorkerBasic) Run() error {
	if len(w.consumers) < 1 {
		xlog.Error("no worker in register")
		return errors.New("no worker in register")
	}
	var ctx context.Context
	ctx, w.CancelFunc = context.WithCancel(context.Background())
	for _, v := range w.consumers {
		go func(c ConsumerFunc) {
			xlog.Info("run worker", xlog.String("x-worker-name", c.GetName()))
			c.Init()
			err := c.Serve(ctx)
			if err != nil {
				xlog.Panic(err.Error(), xlog.FieldErr(err))
				//panic(err)
			}
		}(v)
	}
	return nil
}
func (w *WorkerBasic) Stop() error {
	w.CancelFunc()
	return nil
}
func NewWorker(c ...ConsumerFunc) *WorkerBasic {
	err := InitPool()
	if err != nil {
		xlog.Panic(err.Error(), xlog.FieldErr(err))
		//panic(err)
	}
	return &WorkerBasic{consumers: c}
}
func (w *WorkerBasic) RegisterWorker(c ...ConsumerFunc) {
	w.consumers = append(w.consumers, c...)
}

func InitPool() (topErr error) {
	defer func() {
		if err := recover(); err != nil {
			topErr = fmt.Errorf("panic in InitAntsPool err:%#v", err)
		}
	}()
	var n int
	if r := config.Get("worker.num"); r != nil && r.(int) > 0 {
		if nr, ok := r.(int); ok {
			n = nr
			_, err := xpool.NewAntsPool("workerGlobal", n, xpool.WithAntsLogger())
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("get config worker.num err")
}
