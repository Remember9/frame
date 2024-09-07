package tracing

import (
	"context"
	"esfgit.leju.com/golang/frame/config"
	provider2 "esfgit.leju.com/golang/frame/util/xmiddware/xtrace/provider"
	"esfgit.leju.com/golang/frame/util/xtransport"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Option is tracing option.
type Option func(*options)

type options struct {
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

// WithPropagator with tracer propagator.
func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(opts *options) {
		opts.propagator = propagator
	}
}

// WithTracerProvider with tracer provider.
// Deprecated: use otel.SetTracerProvider(provider) instead.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(opts *options) {
		opts.tracerProvider = provider
	}
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) endpoint.Middleware {
	tracer := NewTracer(trace.SpanKindServer, opts...)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (reply interface{}, err error) {
			if tr, ok := xtransport.FromServerContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				setServerSpan(ctx, span, request)
				defer func() { tracer.End(ctx, span, reply, err) }()
			}
			return next(ctx, request)
		}

	}
}
func Client(tracer *Tracer, opts ...Option) endpoint.Middleware {
	//tracer := NewTracer(trace.SpanKindClient, opts...)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (reply interface{}, err error) {
			if tr, ok := xtransport.FromClientContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				//span.SetName(operationName)
				SetClientSpan(ctx, span, request)
				defer func() { tracer.End(ctx, span, reply, err) }()
			}
			return next(ctx, request)
		}
	}
}

var IsTracing bool

func InitTracing() func() {
	url := "http://localhost:9411/api/v2/spans"
	if r := config.Get("tracing.url"); r != nil {
		url = r.(string)
	}
	if r := config.Get("tracing.isuse"); r != nil {
		if r.(bool) == true {
			IsTracing = true
			shutdown := provider2.InitZipkinTracer(url)
			return shutdown
		}
	}
	return nil
}
