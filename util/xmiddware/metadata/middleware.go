package metadata

import (
	"context"
	"github.com/Remember9/frame/util/xtransport"
	"github.com/go-kit/kit/endpoint"
	"strings"
)

// Option is metadata option.
type Option func(*options)

type options struct {
	prefix []string
	md     Metadata
}

func (o *options) hasPrefix(key string) bool {
	k := strings.ToLower(key)
	for _, prefix := range o.prefix {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// WithConstants with constant metadata key value.
func WithConstants(md Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// WithPropagatedPrefix with propagated key prefix.
func WithPropagatedPrefix(prefix ...string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Server is middleware server-side metadata.
func Server(opts ...Option) endpoint.Middleware {
	options := &options{
		prefix: []string{"x-md-"}, // x-md-global-, x-md-local
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := xtransport.FromServerContext(ctx); ok {
				md := options.md.Clone()
				header := tr.RequestHeader()
				for _, k := range header.Keys() {
					if options.hasPrefix(k) {
						md.Set(k, header.Get(k))
					}
				}
				ctx = NewServerContext(ctx, md)
			}
			return handler(ctx, req)
		}
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) endpoint.Middleware {
	options := &options{
		prefix: []string{"x-md-global-"},
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := xtransport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				// x-md-local-
				for k, v := range options.md {
					header.Set(k, v)
				}
				if md, ok := FromClientContext(ctx); ok {
					for k, v := range md {
						header.Set(k, v)
					}
				}
				// x-md-global-
				if md, ok := FromServerContext(ctx); ok {
					for k, v := range md {
						if options.hasPrefix(k) {
							header.Set(k, v)
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
