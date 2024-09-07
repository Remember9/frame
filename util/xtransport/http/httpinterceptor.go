package http

import (
	"context"
	"github.com/Remember9/frame/util/xtransport"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

func HttpServerFilter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			ctx, cancel = context.WithCancel(req.Context())
			defer cancel()
			pathTemplate := req.URL.Path
			if route := mux.CurrentRoute(req); route != nil {
				pathTemplate, _ = route.GetPathTemplate()
			}
			tr := NewTransport(pathTemplate, pathTemplate, pathTemplate, ToheaderCarrier(req.Header), ToheaderCarrier(w.Header()), req)
			/*tr := &Transport{
				endpoint:     pathTemplate,
				operation:    pathTemplate,
				reqHeader:    headerCarrier(req.Header),
				replyHeader:  headerCarrier(w.Header()),
				request:      req,
				pathTemplate: pathTemplate,
			}*/
			ctx = xtransport.NewServerContext(ctx, tr)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
func HttpClientFilter(mid []endpoint.Middleware) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, req *http.Request) context.Context {
		tr := NewTransport(req.Host, req.URL.Path, req.URL.Path, ToheaderCarrier(req.Header), nil, req)
		ctx = xtransport.NewClientContext(ctx, tr)
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		if len(mid) > 0 {
			mchain := endpoint.Chain(mid[0], mid[1:]...)
			h = mchain(h)
		}
		_, _ = h(ctx, req)
		return ctx
	}
}
