package gorm

import (
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/util/xcast"
	"esfgit.leju.com/golang/frame/xlog"
	"gorm.io/gorm"
	"time"
)

// config options
type Config struct {
	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN         string
	DSNReplicas []string
	// Debug开关
	LogLevel string
	// 最大空闲连接数
	MaxIdleConns int
	// 最大活动连接数
	MaxOpenConns int
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration
	// 慢日志阈值
	SlowThreshold time.Duration
	// 生成 SQL 但不执行
	DryRun bool
	// 跳过默认事务
	SkipDefaultTransaction bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:                    "",
		DSNReplicas:            []string{},
		LogLevel:               "warn",
		MaxIdleConns:           10,
		MaxOpenConns:           100,
		ConnMaxLifetime:        xcast.ToDuration("300s"),
		SlowThreshold:          xcast.ToDuration("500ms"),
		DryRun:                 false,
		SkipDefaultTransaction: true,
	}
}

func Build(name string) *gorm.DB {
	var dbConfig = DefaultConfig()
	if err := config.UnmarshalKey(name, dbConfig); err != nil {
		xlog.Panic("unmarshal key", xlog.String("mod", "gorm"), xlog.FieldErr(err), xlog.String("key", name))
	}
	dbConfig.Name = name

	db, err := OpenMysql(dbConfig)
	if err != nil {
		xlog.Panic("dial nil", xlog.String("mod", "gorm"), xlog.FieldErr(err), xlog.Any("value", dbConfig))
		return db
	}
	if db == nil {
		xlog.Panic("db nil", xlog.String("mod", "gorm"), xlog.Any("value", dbConfig))
		return db
	}
	return db
}
