package grpc

import (
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	"github.com/Remember9/frame/util/xmiddware"
	grpc2 "github.com/Remember9/frame/util/xtransport/grpc"
	"github.com/Remember9/frame/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Config ...
type Config struct {
	Name        string
	Address     string
	Block       bool
	DialTimeout time.Duration
	ReadTimeout time.Duration
	KeepAlive   *keepalive.ClientParameters
	dialOptions []grpc.DialOption

	SlowThreshold time.Duration
	Debug         bool
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		DialTimeout:   time.Second * 3,
		ReadTimeout:   xcast.ToDuration("1s"),
		SlowThreshold: xcast.ToDuration("600ms"),
		Block:         true,
	}
}

func Build(name string) *grpc.ClientConn {
	var grpcClientConfig = DefaultConfig()
	if err := config.UnmarshalKey("client."+name, &grpcClientConfig); err != nil {
		xlog.Panic("client grpc parse config panic", xlog.String("err kind", "unmarshal config err"), xlog.FieldErr(err), xlog.String("key", name), xlog.Any("value", grpcClientConfig))
	}
	grpcClientConfig.Name = name

	// 链路追踪和元数据传递
	middlewares := xmiddware.GetGrpcClinetMiddleware()
	grpcClientConfig.dialOptions = append(grpcClientConfig.dialOptions,
		grpc.WithChainUnaryInterceptor(grpc2.UnaryGrpcClientInterceptor(0, middlewares)))

	return newGRPCClient(grpcClientConfig)
}
