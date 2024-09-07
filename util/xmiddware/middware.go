package xmiddware

import (
	"context"
	"esfgit.leju.com/golang/frame/util/xmiddware/metadata"
	"esfgit.leju.com/golang/frame/util/xmiddware/xtrace/tracing"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

// Middleware 用于服务端

var serverMiddleware []endpoint.Middleware
var clientTracer = tracing.NewTracer(trace.SpanKindClient)
var httpClientMiddleware []endpoint.Middleware
var grpcClientMiddleware []endpoint.Middleware

func GetServerMiddleware() []endpoint.Middleware {
	if serverMiddleware == nil {
		serverMiddleware = append(serverMiddleware, NopFirst(), metadata.Server())
		if tracing.IsTracing {
			serverMiddleware = append(serverMiddleware, tracing.Server())
		}
	}
	return serverMiddleware
}

func GetHttpClientMiddleware() []endpoint.Middleware {
	if httpClientMiddleware == nil {
		serverMiddleware = append(serverMiddleware, NopFirst(), metadata.Server())
		httpClientMiddleware = append(httpClientMiddleware, NopFirst(), metadata.Client())
		if tracing.IsTracing {
			httpClientMiddleware = append(httpClientMiddleware, tracing.Client(clientTracer))
		}
	}
	return httpClientMiddleware
}

func GetGrpcClinetMiddleware() []endpoint.Middleware {
	if grpcClientMiddleware == nil {
		serverMiddleware = append(serverMiddleware, NopFirst(), metadata.Server())
		grpcClientMiddleware = append(grpcClientMiddleware, NopFirst(), metadata.Client())
		if tracing.IsTracing {
			grpcClientMiddleware = append(grpcClientMiddleware, tracing.Client(clientTracer))
		}
	}
	return grpcClientMiddleware
}

func NopFirst() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return e(ctx, request)
		}
	}
}
