package xresty

import "github.com/go-resty/resty/v2"

func NewClient() *resty.Client {
	cnf := getConfig()
	return resty.New().SetLogger(DefaultLogger()).
		SetTimeout(cnf.Timeout).
		SetRetryCount(cnf.RetryCount).
		SetDebug(cnf.Debug)
}
