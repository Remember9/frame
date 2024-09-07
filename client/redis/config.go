package redis

import (
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	"github.com/Remember9/frame/xlog"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	// ClusterMode using clusterClient
	ClusterMode string = "cluster"
	// StubMode using redisClient
	StubMode string = "stub"
)

// Config for redis, contains RedisStubConfig and RedisClusterConfig
type Config struct {
	// Addrs 实例配置地址
	Addrs []string
	// Mode Redis模式 cluster|stub
	Mode string
	// Password 密码
	Password string
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int
	// PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	PoolSize int
	// MaxRetries 网络相关的错误最大重试次数 默认8次
	MaxRetries int
	// MinIdleConns 最小空闲连接数
	MinIdleConns int
	// DialTimeout 拨超时时间
	DialTimeout time.Duration
	// ReadTimeout 读超时 默认3s
	ReadTimeout time.Duration
	// WriteTimeout 读超时 默认3s
	WriteTimeout time.Duration
	// IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	IdleTimeout time.Duration
	// Debug开关
	Debug bool
	// ReadOnly 集群模式 在从属节点上启用读模式
	ReadOnly bool
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration
}

// DefaultRedisConfig default config ...
func DefaultRedisConfig() *Config {
	return &Config{
		DB:            0,
		PoolSize:      10,
		MaxRetries:    3,
		MinIdleConns:  100,
		DialTimeout:   xcast.ToDuration("1s"),
		ReadTimeout:   xcast.ToDuration("1s"),
		WriteTimeout:  xcast.ToDuration("1s"),
		IdleTimeout:   xcast.ToDuration("60s"),
		ReadOnly:      false,
		Debug:         false,
		SlowThreshold: xcast.ToDuration("250ms"),
	}
}

func Build(name string) *Redis {
	var redisConfig = DefaultRedisConfig()
	if err := config.UnmarshalKey(name, &redisConfig); err != nil {
		xlog.Panic("unmarshal redisConfig",
			xlog.String("key", name),
			xlog.Any("redisConfig", redisConfig),
			xlog.String("error", err.Error()))
	}
	count := len(redisConfig.Addrs)
	if count < 1 {
		xlog.Panic("no address in redis config", xlog.Any("config", redisConfig))
	}
	if len(redisConfig.Mode) == 0 {
		redisConfig.Mode = StubMode
		if count > 1 {
			redisConfig.Mode = ClusterMode
		}
	}
	redis.SetLogger(&WrapLogger{})
	var client redis.Cmdable
	switch redisConfig.Mode {
	case ClusterMode:
		if count == 1 {
			xlog.Warn("redis config has only 1 address but with cluster mode")
		}
		client = redisConfig.buildCluster()
	case StubMode:
		if count > 1 {
			xlog.Warn("redis config has more than 1 address but with stub mode")
		}
		client = redisConfig.buildStub()
	default:
		xlog.Panic("redis mode must be one of (stub, cluster)")
	}
	return &Redis{
		Config: redisConfig,
		Client: client,
	}
}

func (config Config) buildStub() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.Addrs[0],
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
}

func (config Config) buildCluster() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		MaxRedirects: config.MaxRetries,
		ReadOnly:     config.ReadOnly,
		Password:     config.Password,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
}
