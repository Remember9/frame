package grpc

import (
	"github.com/Remember9/frame/util/xtransport"
	"google.golang.org/grpc/metadata"
)

var _ xtransport.Transporter = &Transport{}

// Transport is a gRPC transport.
type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
	// filters     []selector.NodeFilter
}

func NewTransport(endpoint, operation string, reqHeader, replyHeader headerCarrier) *Transport {
	return &Transport{
		endpoint:    endpoint,
		operation:   operation,
		reqHeader:   reqHeader,
		replyHeader: replyHeader,
	}
}

// Kind returns the transport kind.
func (tr *Transport) Kind() xtransport.Kind {
	return xtransport.KindGRPC
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() xtransport.Header {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() xtransport.Header {
	return tr.replyHeader
}

// Filters returns the client select filters.
/*func (tr *Transport) NodeFilters() []selector.NodeFilter {
	return tr.filters
}*/

type headerCarrier metadata.MD

func ToheaderCarrier(m metadata.MD) headerCarrier {
	if m == nil {
		return headerCarrier(metadata.MD{})
	}
	return headerCarrier(m)
}

// Get returns the value associated with the passed key.
func (mc headerCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc headerCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}
