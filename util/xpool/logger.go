package xpool

import (
	"github.com/Remember9/frame/xlog"
	"github.com/panjf2000/ants/v2"
)

type AntsLogger struct {
}

func (wl *AntsLogger) Printf(format string, args ...interface{}) {
	// 目前ans记录日志时都是panic时记录，所以暂时用errorf
	xlog.Errorf(format, args)
}

func WithAntsLogger() ants.Option {
	return ants.WithLogger(&AntsLogger{})
}
