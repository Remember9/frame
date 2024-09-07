package grpc

import (
	"context"
	"encoding/json"
	"github.com/Remember9/frame/util/xerrors"
	"github.com/Remember9/frame/util/xstring"
	"github.com/Remember9/frame/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

// timeoutUnaryClientInterceptor gRPC客户端超时拦截器
func timeoutUnaryClientInterceptor(timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		// 若无自定义超时设置，默认设置超时
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		du := time.Since(now)
		remoteIP := "unknown"
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			remoteIP = remote.Addr.String()
		}

		if slowThreshold > time.Duration(0) && du > slowThreshold {
			xlog.Error("slow",
				xlog.String("err", "grpc unary slow command"),
				xlog.String("method", method),
				xlog.String("name", cc.Target()),
				xlog.Duration("cost", du),
				xlog.String("addr", remoteIP),
			)
		}

		return err
	}
}

// loggerUnaryClientInterceptor gRPC客户端日志中间件
func loggerUnaryClientInterceptor(name string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		spbStatus := xerrors.ToCodeError(err)
		if err != nil {
			xlog.Error(
				"Access",
				xlog.String("type", "unary"),
				xlog.Int64("code", spbStatus.Code()),
				xlog.String("err", spbStatus.Error()),
				xlog.String("name", name),
				xlog.String("method", method),
				xlog.Duration("cost", time.Since(beg)),
				xlog.Any("req", json.RawMessage(xstring.Json(req))),
				xlog.Any("reply", json.RawMessage(xstring.Json(reply))),
			)
			return err
		} else {
			xlog.Debug(
				"Access",
				xlog.String("type", "unary"),
				xlog.Int64("code", 0),
				xlog.String("name", name),
				xlog.String("method", method),
				xlog.Duration("cost", time.Since(beg)),
				xlog.Any("req", json.RawMessage(xstring.Json(req))),
				xlog.Any("reply", json.RawMessage(xstring.Json(reply))),
			)

		}

		return nil
	}
}
