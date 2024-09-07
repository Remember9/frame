package http

import (
	"bytes"
	"context"
	"encoding/json"
	"esfgit.leju.com/golang/frame/util/xencoding"
	"esfgit.leju.com/golang/frame/util/xmiddware"
	http2 "esfgit.leju.com/golang/frame/util/xtransport/http"
	httputil2 "esfgit.leju.com/golang/frame/util/xtransport/httputil"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewHttpClinet(method string, tgt *url.URL, replyType interface{}, options ...httptransport.ClientOption) endpoint.Endpoint {

	options = append(options, httptransport.ClientBefore(http2.HttpClientFilter(xmiddware.GetHttpClientMiddleware())))

	e := httptransport.NewClient(method, tgt, httptransport.EncodeJSONRequest, DefaultResponseDecoder(replyType), options...).Endpoint()
	//e=tracing.Client(e)
	return e
}

func DefaultRequestEncoder(tgt *url.URL) func(context.Context, *http.Request, interface{}) error {
	return func(ctx context.Context, req *http.Request, request interface{}) error {
		var buf bytes.Buffer
		req.URL = tgt
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(&buf)
		return nil
	}
}

// DefaultResponseDecoder is an HTTP response decoder.
func DefaultResponseDecoder(v interface{}) func(context.Context, *http.Response) (interface{}, error) {

	return func(ctx context.Context, res *http.Response) (interface{}, error) {
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if err := CodecForResponse(res).Unmarshal(data, v); err != nil {
			return nil, err
		}
		return v, nil
	}
}

// CodecForResponse get encoding.Codec via http.Response
func CodecForResponse(r *http.Response) xencoding.Codec {
	codec := xencoding.GetCodec(httputil2.ContentSubtype(r.Header.Get("Content-Type")))
	if codec != nil {
		return codec
	}
	return xencoding.GetCodec("json")
}
