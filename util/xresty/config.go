package xresty

import (
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/xlog"
	"sync"
	"time"
)

const key = "resty"

var _cnf *xConfig // 不要直接调用，使用 getConfig 调用
var cnfMut sync.Mutex

type xConfig struct {
	Debug      bool
	RetryCount int
	Timeout    time.Duration // A Timeout of zero means no timeout.
}

func defaultConfig() *xConfig {
	return &xConfig{
		Debug:      true,
		RetryCount: 1,
		Timeout:    20 * time.Second,
	}
}

func getConfig() *xConfig {
	if _cnf == nil {
		cnfMut.Lock()
		defer cnfMut.Unlock()
		if _cnf == nil {
			build()
		}
	}
	return _cnf
}

func build() {
	_cnf = defaultConfig()
	err := config.UnmarshalKey(key, _cnf)
	if err != nil {
		xlog.Panic("unmarshal restyConfig",
			xlog.String("key", key),
			xlog.Any("restyConfig", _cnf),
			xlog.String("error", err.Error()))
	}
}
