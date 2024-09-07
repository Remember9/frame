package redis

import (
	"github.com/Remember9/frame/xlog"
	"golang.org/x/net/context"
)

type WrapLogger struct{}

func (l *WrapLogger) Printf(_ context.Context, format string, v ...interface{}) {
	xlog.Debugf("redis:"+format, v...)
}
