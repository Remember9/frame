package grpc

import (
	"context"
	"esfgit.leju.com/golang/frame/xlog"
	"google.golang.org/grpc"
	"time"
)

func newGRPCClient(c *Config) *grpc.ClientConn {
	var ctx = context.Background()
	c.dialOptions = append(c.dialOptions,
		grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(c.ReadTimeout, c.SlowThreshold)),
		grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(c.Name)),
	)
	if c.Block {
		if c.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.DialTimeout)
			defer cancel()
		}
		c.dialOptions = append(c.dialOptions, grpc.WithBlock())
	}
	cc, err := grpc.DialContext(ctx, c.Address, c.dialOptions...)

	if err != nil {
		xlog.Error("dial grpc server",
			xlog.String("mod", "client.grpc"),
			xlog.String("addr", c.Address),
			xlog.String("errKind", "request err"),
			xlog.FieldErr(err))
	}
	xlog.Info("start grpc client",
		xlog.String("mod", "client.grpc"),
		xlog.String("addr", c.Address))
	return cc
}
