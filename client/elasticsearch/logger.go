package elasticsearch

import (
	"github.com/Remember9/frame/xlog"
)

type WrapErrorLogger struct{}

func (logger WrapErrorLogger) Printf(format string, vars ...interface{}) {
	xlog.Errorf(format, vars...)
}

type WrapInfoLogger struct{}

func (logger WrapInfoLogger) Printf(format string, vars ...interface{}) {
	xlog.Infof(format, vars...)
}

type WrapTraceLogger struct{}

func (logger WrapTraceLogger) Printf(format string, vars ...interface{}) {
	xlog.Debugf(format, vars...)
}
