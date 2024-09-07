package xlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Field = zap.Field
	Level = zapcore.Level
)

var (
	Bool        = zap.Bool
	Bools       = zap.Bools
	String      = zap.String
	Strings     = zap.Strings
	ByteString  = zap.ByteString
	ByteStrings = zap.ByteStrings
	Any         = zap.Any
	Array       = zap.Array
	Int         = zap.Int
	Ints        = zap.Ints
	Uint        = zap.Uint
	Uints       = zap.Uints
	Int64       = zap.Int64
	Int64s      = zap.Int64s
	Uint64      = zap.Uint64
	Uint64s     = zap.Uint64s
	Int32       = zap.Int32
	Int32s      = zap.Int32s
	Uint32      = zap.Uint32
	Uint32s     = zap.Uint32s
	Int16       = zap.Int16
	Int16s      = zap.Int16s
	Uint16      = zap.Uint16
	Uint16s     = zap.Uint16s
	Int8        = zap.Int8
	Int8s       = zap.Int8s
	Uint8       = zap.Uint8
	Uint8s      = zap.Uint8s
	Float64     = zap.Float64
	Float64s    = zap.Float64s
	Float32     = zap.Float32
	Float32s    = zap.Float32s
	Duration    = zap.Duration
	Durations   = zap.Durations
	Time        = zap.Time
	Times       = zap.Times
	Object      = zap.Object
	FieldErr    = zap.Error
	FieldErrs   = zap.Errors
)

func Debug(msg string, field ...Field) {
	logger.Debug(msg, field...)
}

func Info(msg string, field ...Field) {
	logger.Info(msg, field...)
}

func Warn(msg string, field ...Field) {
	logger.Warn(msg, field...)
}

func Error(msg string, field ...Field) {
	logger.Error(msg, field...)
}

func Panic(msg string, field ...Field) {
	logger.Panic(msg, field...)
}

func Debugf(template string, args ...interface{}) {
	logger.Sugar().Debugf(sprintf(template, args...))
}

func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(sprintf(template, args...))
}

func Warnf(template string, args ...interface{}) {
	logger.Sugar().Warnf(sprintf(template, args...))
}

func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(sprintf(template, args...))
}

func Panicf(template string, args ...interface{}) {
	logger.Sugar().Panicf(sprintf(template, args...))
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Errorw(msg, keysAndValues...)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}
