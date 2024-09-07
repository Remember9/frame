package gorm

import (
	"context"
	"github.com/Remember9/frame/xlog"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func NewLogger(level gormlogger.LogLevel, slowThreshold time.Duration) gormlogger.Interface {
	return &Logger{
		LogLevel:      level,
		SlowThreshold: slowThreshold,
	}
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	logger := *l
	logger.LogLevel = level
	return &logger
}

func (l Logger) Info(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	xlog.Infof(str, args...)
}

func (l Logger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	xlog.Warnf(str, args...)
}

func (l Logger) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	xlog.Errorf(str, args...)
}

func (l Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	cost := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		sql, rows := fc()
		xlog.Error("gorm trace Error", xlog.FieldErr(err), xlog.Duration("cost", cost), xlog.Int64("rows", rows), xlog.String("sql", sql))
	case l.SlowThreshold != 0 && cost > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		xlog.Warn("gorm trace Warn", xlog.Duration("cost", cost), xlog.Duration("slowThreshold", l.SlowThreshold), xlog.Int64("rows", rows), xlog.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		xlog.Debug("gorm trace Info", xlog.Duration("cost", cost), xlog.Int64("rows", rows), xlog.String("sql", sql))
	}
}

func UnmarshalText(text string) gormlogger.LogLevel {
	switch text {
	case "silent", "SILENT":
		return gormlogger.Silent
	case "error", "ERROR":
		return gormlogger.Error
	case "warn", "WARN", "":
		return gormlogger.Warn
	case "info", "INFO":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}
