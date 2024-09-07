package xresty

import (
	"esfgit.leju.com/golang/frame/xlog"
	"github.com/go-resty/resty/v2"
)

var internalLogger resty.Logger = logger{}

func DefaultLogger() resty.Logger {
	return internalLogger
}

type logger struct {
}

func (l logger) Errorf(format string, v ...interface{}) {
	xlog.Errorf(format, v...)
}
func (l logger) Warnf(format string, v ...interface{}) {
	xlog.Warnf(format, v...)
}
func (l logger) Debugf(format string, v ...interface{}) {
	xlog.Debugf(format, v...)
}
