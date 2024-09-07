package xmux

import (
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	"github.com/Remember9/frame/util/xtransport/http"
	"github.com/Remember9/frame/xlog"
	"time"
)

const ModName = "server.xmux"

// Config HTTP config
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	SlowQueryThresholdInMilli time.Duration
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Addr:                      "127.0.0.1:8080",
		ReadTimeout:               xcast.ToDuration("1s"),
		WriteTimeout:              xcast.ToDuration("1s"),
		SlowQueryThresholdInMilli: xcast.ToDuration("500ms"),
	}
}

func Build() *Server {
	var httpServerConfig = DefaultConfig()
	if err := config.UnmarshalKey("server.http", &httpServerConfig); err != nil {
		xlog.Panic("http server parse config panic", xlog.String("err kind", "unmarshal config err"), xlog.FieldErr(err), xlog.String("key", ModName), xlog.Any("value", httpServerConfig))
	}

	server := newServer(httpServerConfig)

	mw := muxMiddleware{}
	mw.Populate(httpServerConfig.SlowQueryThresholdInMilli)

	server.Router.Use(mw.recoverMiddleware)

	server.Router.Use(http.HttpServerFilter())
	return server
}
