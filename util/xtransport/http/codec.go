package http

import (
	"errors"
	encoding "esfgit.leju.com/golang/frame/util/xencoding"
	httputil2 "esfgit.leju.com/golang/frame/util/xtransport/httputil"
	"io"
	"net/http"
	/*_ "esfgit.leju.com/golang/frame/util/xencoding/form"
	_ "esfgit.leju.com/golang/frame/util/xencoding/json"
	_ "esfgit.leju.com/golang/frame/util/xencoding/proto"
	_ "esfgit.leju.com/golang/frame/util/xencoding/xml"
	_ "esfgit.leju.com/golang/frame/util/xencoding/yaml"*/)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// DefaultRequestDecoder decodes the request body to object.
func defaultRequestDecoder(r *http.Request, v interface{}) error {
	codec, ok := CodecForRequest(r, "Content-Type")
	/*
		a:=r.Header["Content-Type"][0]
		c := encoding.GetCodec(httputil.ContentSubtype(a))
		return errors.New("CODEC:"+fmt.Sprintf("%#v , ######,%#V,########,%#V",codec,r.Header["Content-Type"][0],c))*/

	if !ok {
		return errors.New("CODEC:" + r.Header.Get("Content-Type"))
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("CODEC:" + err.Error())
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.New("CODEC:" + err.Error())
	}
	return nil
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	codec, _ := CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil2.ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the HTTP response.

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(httputil2.ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}
