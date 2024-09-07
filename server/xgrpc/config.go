package xgrpc

import (
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	grpc2 "github.com/Remember9/frame/util/xtransport/grpc"
	"github.com/Remember9/frame/xlog"
	"google.golang.org/grpc"
	"time"
)

const ModName = "server.xgrpc"

// Config ...
type Config struct {
	Addr                      string
	Network                   string
	SlowQueryThresholdInMilli time.Duration
	serverOptions             []grpc.ServerOption
	streamInterceptors        []grpc.StreamServerInterceptor
	unaryInterceptors         []grpc.UnaryServerInterceptor
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                   "tcp4",
		Addr:                      "127.0.0.1:9090",
		SlowQueryThresholdInMilli: xcast.ToDuration("500ms"),
		serverOptions:             []grpc.ServerOption{},
		streamInterceptors:        []grpc.StreamServerInterceptor{},
		unaryInterceptors:         []grpc.UnaryServerInterceptor{},
	}
}
func Build() *Server {
	var grpcServerConfig = DefaultConfig()
	if err := config.UnmarshalKey("server.grpc", &grpcServerConfig); err != nil {
		xlog.Panic("grpc server parse config panic",
			xlog.String("err kind", "unmarshal config err"),
			xlog.FieldErr(err), xlog.String("key", ModName),
			xlog.Any("value", grpcServerConfig),
		)
	}
	// 消息数据报
	grpcServerConfig.unaryInterceptors = append(grpcServerConfig.unaryInterceptors, grpc2.UnaryGrpcServerInterceptor())

	server := newServer(grpcServerConfig)

	return server
}
