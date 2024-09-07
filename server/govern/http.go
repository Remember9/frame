package govern

import (
	"encoding/json"
	"github.com/arl/statsviz"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"
	"strings"
)

var (
	DefaultServeMux = http.NewServeMux()
	routes          []string
)

func init() {
	// 获取全部治理路由
	HandleFunc("/routes", func(resp http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(resp).Encode(routes)
	})

	// pprof
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)
	Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	Handle("/debug/pprof/block", pprof.Handler("block"))
	Handle("/debug/pprof/heap", pprof.Handler("heap"))
	Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	//
	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		envs := os.Environ()
		reply := map[string]string{}
		for k, v := range envs {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) != 2 {
				reply["zzzzzz-"+string(k)] = v
			} else {
				reply[parts[0]] = parts[1]
			}
		}
		_ = jsoniter.NewEncoder(w).Encode(reply)
	})
	HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	if info, ok := debug.ReadBuildInfo(); ok {
		HandleFunc("/mod/info", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(info)
		})
	}

	// statsviz
	HandleFunc("/debug/statsviz/", statsviz.Index)
	HandleFunc("/debug/statsviz/ws", statsviz.Ws)

}

func HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}

func Handle(pattern string, handler http.Handler) {
	DefaultServeMux.Handle(pattern, handler)
	routes = append(routes, pattern)
}
