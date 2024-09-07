package http

import (
	"context"
	"github.com/Remember9/frame/util/xtransport"
	"net/http"
)

var _ Transporter = &Transport{}

// Transporter is http Transporter
type Transporter interface {
	xtransport.Transporter
	Request() *http.Request
	PathTemplate() string
}

// Transport is an HTTP transport.
type Transport struct {
	endpoint     string
	operation    string
	reqHeader    headerCarrier
	replyHeader  headerCarrier
	request      *http.Request
	pathTemplate string
}

func NewTransport(endpoint, operation, pathTemplate string, reqHeader, replyHeader headerCarrier, request *http.Request) *Transport {
	return &Transport{
		endpoint:     endpoint,
		operation:    operation,
		reqHeader:    reqHeader,
		replyHeader:  replyHeader,
		request:      request,
		pathTemplate: pathTemplate,
	}
}

// Kind returns the transport kind.
func (tr *Transport) Kind() xtransport.Kind {
	return xtransport.KindHTTP
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Request returns the HTTP request.
func (tr *Transport) Request() *http.Request {
	return tr.request
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() xtransport.Header {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() xtransport.Header {
	return tr.replyHeader
}

// PathTemplate returns the http path template.
func (tr *Transport) PathTemplate() string {
	return tr.pathTemplate
}

// SetOperation sets the transport operation.
func SetOperation(ctx context.Context, op string) {
	if tr, ok := xtransport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			tr.operation = op
		}
	}
}

type headerCarrier http.Header

func ToheaderCarrier(h http.Header) headerCarrier {
	return headerCarrier(h)
}

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}
