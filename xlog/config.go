package xlog

import (
	"github.com/Remember9/frame/config"
)

var logger *esflogger

type Config struct {
	Prefix         string
	Level          string
	PathDir        string
	Development    bool
	DisableCaller  bool
	ConsoleEncoder bool
	MaxSize        int
	MaxBackups     int
	MaxAge         int
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		PathDir:        "./logs/",
		Prefix:         "frame-srv",
		Level:          "info",
		MaxSize:        100,
		MaxBackups:     60,
		MaxAge:         30,
		Development:    true,
		DisableCaller:  true,
		ConsoleEncoder: false,
	}
}

// Build ...
func Build() error {
	var logConfig = DefaultConfig()
	if err := config.UnmarshalKey("log", &logConfig); err != nil {
		panic(err)
	}

	op := &options{
		LogFileDir:     logConfig.PathDir,
		Prefix:         logConfig.Prefix,
		MaxSize:        logConfig.MaxSize,
		MaxBackups:     logConfig.MaxBackups,
		MaxAge:         logConfig.MaxAge,
		Development:    logConfig.Development,
		PanicFileName:  "panic.log",
		ErrorFileName:  "error.log",
		WarnFileName:   "warn.log",
		InfoFileName:   "info.log",
		DebugFileName:  "debug.log",
		DisableCaller:  logConfig.DisableCaller,
		ConsoleEncoder: logConfig.ConsoleEncoder,
	}
	_ = op.Level.UnmarshalText([]byte(logConfig.Level))

	logger = newLogger(op)

	return nil
}
