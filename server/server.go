package server

import (
	"context"
	"fmt"
	"strings"
)

// ServiceInfo ...
type ServiceInfo struct {
	Name    string
	Scheme  string
	Address string
}

// Label ...
func (si ServiceInfo) Label() string {
	address := si.Address
	if strings.HasPrefix(address, ":") {
		address = "127.0.0.1" + address
	}
	return fmt.Sprintf("%s://%s", si.Scheme, address)
}

// Server ...
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}
