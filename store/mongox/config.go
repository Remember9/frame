package mongox

import (
	"context"
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/xlog"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

// Config for initial mongo instance
type Config struct {
	// URI example: mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	// URI Reference: https://docs.mongodb.com/manual/reference/connection-string/
	Uri      string `json:"uri"`
	Database string `json:"database"`
	Coll     string `json:"coll"`
	// ConnectTimeoutMS specifies a timeout that is used for creating connections to the server.
	//	If set to 0, no timeout will be used.
	//	The default is 30 seconds.
	ConnectTimeoutMS *int64 `json:"connectTimeoutMS"`
	// MaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
	// If this is 0, it will be set to math.MaxInt64,
	// The default is 100.
	MaxPoolSize *uint64 `json:"maxPoolSize"`
	// MinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
	// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
	// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
	MinPoolSize *uint64 `json:"minPoolSize"`
	// SocketTimeoutMS specifies how long the driver will wait for a socket read or write to return before returning a
	// network error. If this is 0 meaning no timeout is used and socket operations can block indefinitely.
	// The default is 300,000 ms.
	SocketTimeoutMS *int64 `json:"socketTimeoutMS"`
}

// 默认配置
func DefaultConfig() *qmgo.Config {
	var cTimeout int64 = 30
	var sTimeout int64 = 300000
	var maxPoolSize uint64 = 100
	var minPoolSize uint64 = 0
	return &qmgo.Config{
		Uri:              "",
		ConnectTimeoutMS: &cTimeout,
		SocketTimeoutMS:  &sTimeout,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
	}
}

func rawConfig(name string) *qmgo.Config {
	var mongoConfig = DefaultConfig()
	if err := config.UnmarshalKey(name, mongoConfig); err != nil {
		xlog.Panic("unmarshal key", xlog.String("mod", "mongo"), xlog.FieldErr(err), xlog.String("key", name))
	}
	return mongoConfig
}

func Build(ctx context.Context, name string, o ...options.ClientOptions) *Client {
	mongoConfig := rawConfig(name)
	client, err := qmgo.NewClient(ctx, mongoConfig, o...)
	if err != nil {
		xlog.Panic("connect mongo", xlog.String("mod", "mongo"), xlog.FieldErr(err), xlog.Any("config", mongoConfig))
	}

	return &Client{client, mongoConfig}
}

func BuildCli(ctx context.Context, name string, o ...options.ClientOptions) *Cli {
	mongoConfig := rawConfig(name)
	client, err := qmgo.Open(ctx, mongoConfig, o...)
	if err != nil {
		xlog.Panic("connect mongo cli", xlog.String("mod", "mongo"), xlog.FieldErr(err), xlog.Any("config", mongoConfig))
	}

	return &Cli{client, mongoConfig}
}
