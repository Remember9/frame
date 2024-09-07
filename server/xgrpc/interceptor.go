package xgrpc

import (
	"context"
	"fmt"
	"github.com/Remember9/frame/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"runtime"
	"strings"
	"time"
)

func defaultStreamServerInterceptor(slowQueryThresholdInMilli time.Duration) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]xlog.Field, 0, 8)
		var event = "normal"
		defer func() {
			cast := time.Since(beg)
			if slowQueryThresholdInMilli != 0 && cast > slowQueryThresholdInMilli {
				event = "slow"
			}

			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, xlog.ByteString("stack", stack))
				event = "recover"
			}

			fields = append(fields,
				xlog.Any("grpc interceptor type", "stream"),
				xlog.String("method", info.FullMethod),
				xlog.Duration("cost", cast),
				xlog.String("event", event),
			)

			for key, val := range getPeer(stream.Context()) {
				fields = append(fields, xlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, xlog.String("err", err.Error()))
				xlog.Error("Access", fields...)
				return
			}
			xlog.Debug("Access", fields...)
		}()
		return handler(srv, stream)
	}
}

func defaultUnaryServerInterceptor(slowQueryThresholdInMilli time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var beg = time.Now()
		var fields = make([]xlog.Field, 0, 8)
		var event = "normal"
		defer func() {
			cast := time.Since(beg)
			if slowQueryThresholdInMilli != 0 && cast > slowQueryThresholdInMilli {
				event = "slow"
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, xlog.ByteString("stack", stack))
				event = "recover"
			}

			fields = append(fields,
				xlog.Any("grpc interceptor type", "unary"),
				xlog.String("method", info.FullMethod),
				xlog.Duration("cost", cast),
				xlog.String("event", event),
			)

			for key, val := range getPeer(ctx) {
				fields = append(fields, xlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, xlog.String("err", err.Error()))
				xlog.Error("Access", fields...)
				return
			}
			xlog.Debug("Access", fields...)
		}()
		return handler(ctx, req)
	}
}

func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

func getPeer(ctx context.Context) map[string]string {
	var peerMeta = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			ip, err := getClientIP(ctx)
			if err == nil {
				clientIP = ip
			}
		}
		peerMeta["clientIP"] = clientIP
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta

}
