package govern

import (
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/xlog"
)

// ModName ..
const ModName = "server.govern"

// Config ...
type Config struct {
	Addr string
}

func DefaultConfig() *Config {
	return &Config{
		Addr: "127.0.0.1:8090",
	}
}

func Build() *Server {
	var governServerConfig = DefaultConfig()
	if err := config.UnmarshalKey("server.govern", &governServerConfig); err != nil {
		xlog.Panic("govern server parse config panic", xlog.String("err kind", "unmarshal config err"), xlog.FieldErr(err), xlog.String("key", ModName), xlog.Any("value", governServerConfig))
	}
	return newServer(governServerConfig)
}
