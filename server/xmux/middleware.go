package xmux

import (
	"fmt"
	"esfgit.leju.com/golang/frame/util/xnet"
	"esfgit.leju.com/golang/frame/xlog"
	"net/http"
	"runtime"
	"time"
)

type muxMiddleware struct {
	slowQueryThresholdInMilli time.Duration
}

func (mw *muxMiddleware) Populate(slowQueryThresholdInMilli time.Duration) {
	mw.slowQueryThresholdInMilli = slowQueryThresholdInMilli
}

func (mw *muxMiddleware) recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var beg = time.Now()
		var fields = make([]xlog.Field, 0, 8)

		defer func() {
			sinceTime := time.Since(beg)
			fields = append(fields, xlog.Duration("cost", sinceTime))
			var err error
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, true)
				fields = append(fields, xlog.ByteString("stack", stack[:length]))
			}
			fields = append(fields,
				xlog.String("method", r.Method),
				xlog.String("host", r.Host),
				xlog.String("path", r.URL.Path),
				xlog.String("query", r.URL.RawQuery),
				xlog.String("ip", xnet.ClientIP(r)),
			)

			cost := time.Since(beg)
			if mw.slowQueryThresholdInMilli != 0 && cost > mw.slowQueryThresholdInMilli {
				fields = append(fields, xlog.Duration("slow", cost))
			}

			if err != nil {
				fields = append(fields, xlog.FieldErr(err))
				xlog.Error("Access", fields...)
				return
			}
			xlog.Debug("Access", fields...)

			return
		}()

		h.ServeHTTP(w, r)
	})
}
