package grpc

import (
	"context"
	"esfgit.leju.com/golang/frame/util/xtransport"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
	"time"
)

func UnaryGrpcServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := grpcmd.FromIncomingContext(ctx)
		replyHeader := grpcmd.MD{}
		tr := NewTransport(info.FullMethod, info.FullMethod, ToheaderCarrier(md), ToheaderCarrier(replyHeader))
		ctx = xtransport.NewServerContext(ctx, tr)
		/*ctx = transport.NewServerContext(ctx, &Transport{
			endpoint:    info.FullMethod,
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		})*/
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}

		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

func UnaryGrpcClientInterceptor(timeout time.Duration, mid []endpoint.Middleware) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tr := NewTransport(cc.Target(), method, ToheaderCarrier(nil), nil)
		/*ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:  cc.Target(),
			operation: method,
			reqHeader: headerCarrier{},
		})*/
		ctx = xtransport.NewClientContext(ctx, tr)
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := xtransport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				keyvals := make([]string, 0, len(keys))
				for _, k := range keys {
					keyvals = append(keyvals, k, header.Get(k))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(mid) > 0 {
			mchain := endpoint.Chain(mid[0], mid[1:]...)
			h = mchain(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
