package xresty

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

var FailedRespMiddleware = func(client *resty.Client, resp *resty.Response) error {
	if !resp.IsSuccess() {
		return errors.New(resp.Status())
	}
	return nil
}

var EmptyRespMiddleware = func(client *resty.Client, resp *resty.Response) error {
	if len(resp.Body()) == 0 {
		return errors.New("响应体为空")
	}
	return nil
}
