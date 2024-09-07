package elasticsearch

import (
	"esfgit.leju.com/golang/frame/config"
	"esfgit.leju.com/golang/frame/util/xcast"
	"esfgit.leju.com/golang/frame/xlog"
	"time"
)

// Config ...
type Config struct {
	URL []string
	// supporting v7. Default to v7 if empty.
	Version string
	// optional username to communicate with ElasticSearch
	Username string
	// optional password to communicate with ElasticSearch
	Password string
	// optional to disable sniff, according to issues on Github,
	// Sniff could cause issue like "no Elasticsearch node available"
	DisableSniff bool
	// Gzip
	Compress bool
	// Debug
	Debug bool
	// httpClient
	TransportMaxIdleConns        int
	TransportMaxIdleConnsPerHost int
	TransportTimeout             time.Duration
	Timeout                      time.Duration
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		URL:                          []string{"http://127.0.0.1:9200"},
		TransportMaxIdleConns:        50,
		TransportMaxIdleConnsPerHost: 10,
		TransportTimeout:             xcast.ToDuration("1s"),
		Timeout:                      xcast.ToDuration("1s"),
		DisableSniff:                 true,
		Compress:                     true,
		Debug:                        true,
	}
}

func Build(name string) (GenericClient, error) {
	var esClientConfig = DefaultConfig()
	if err := config.UnmarshalKey(name, &esClientConfig); err != nil {
		xlog.Panic("client elasticsearch parse config panic", xlog.String("err kind", "unmarshal config err"), xlog.FieldErr(err), xlog.String("key", name), xlog.Any("value", esClientConfig))
	}
	return NewGenericClient(esClientConfig)
}
