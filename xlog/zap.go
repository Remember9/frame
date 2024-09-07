package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	l                                       *esflogger
	sp                                      = string(filepath.Separator)
	panicWs, errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                          = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                          = zapcore.Lock(os.Stderr)
)

type (
	options struct {
		LogFileDir     string //文件保存地方
		Prefix         string //日志文件前缀
		PanicFileName  string
		ErrorFileName  string
		WarnFileName   string
		InfoFileName   string
		DebugFileName  string
		Level          zapcore.Level //日志等级
		MaxSize        int           //日志文件小大（M）
		MaxBackups     int           // 最多存在多少个切片文件
		MaxAge         int           //保存的最大天数
		Development    bool          //是否是开发模式
		DisableCaller  bool          //是否禁用Caller
		ConsoleEncoder bool          // encoder为标准控制台
		zap.Config
	}
	esflogger struct {
		*zap.Logger
		sync.RWMutex
		Opts      *options `json:"opts"`
		zapConfig zap.Config
		inited    bool
	}
)

func newLogger(cf *options) *esflogger {
	l = &esflogger{}
	l.Lock()
	defer l.Unlock()
	if l.inited {
		l.Info("logger Inited")
		return nil
	}
	l.Opts = cf
	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	l.zapConfig.DisableCaller = l.Opts.DisableCaller
	l.zapConfig.Level.SetLevel(l.Opts.Level)
	l.init()
	l.inited = true
	l.Info("init logger")
	return l
}

func (l *esflogger) init() {
	l.setSyncers()
	var err error
	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

func (l *esflogger) setSyncers() {
	hostName, _ := os.Hostname()
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.Prefix + "-" + hostName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   false,
			LocalTime:  true,
		})
	}
	panicWs = f(l.Opts.PanicFileName)
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}

func (l *esflogger) cores() zap.Option {
	customEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	if l.Opts.ConsoleEncoder {
		l.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		customEncoder = zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
	}
	panicPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.PanicLevel && zapcore.PanicLevel-l.zapConfig.Level.Level() > -1
	})
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	var cores []zapcore.Core

	if l.Opts.Development {
		cores = []zapcore.Core{
			zapcore.NewCore(customEncoder, errorConsoleWS, panicPriority),
			zapcore.NewCore(customEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(customEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(customEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(customEncoder, debugConsoleWS, debugPriority),
		}
	} else {
		cores = []zapcore.Core{
			zapcore.NewCore(customEncoder, panicWs, panicPriority),
			zapcore.NewCore(customEncoder, errWS, errPriority),
			zapcore.NewCore(customEncoder, warnWS, warnPriority),
			zapcore.NewCore(customEncoder, infoWS, infoPriority),
			zapcore.NewCore(customEncoder, debugWS, debugPriority),
		}
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}
