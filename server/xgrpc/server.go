package xgrpc

import (
	"context"
	"github.com/Remember9/frame/server"
	"github.com/Remember9/frame/xlog"
	"google.golang.org/grpc"
	"net"
)

// Server ...
type Server struct {
	*grpc.Server
	listener net.Listener
	config   *Config
}

func newServer(config *Config) *Server {
	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{defaultStreamServerInterceptor(config.SlowQueryThresholdInMilli)},
		config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor(config.SlowQueryThresholdInMilli)},
		config.unaryInterceptors...,
	)

	config.serverOptions = append(config.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	newServer := grpc.NewServer(config.serverOptions...)
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		xlog.Panic("new "+ModName+" server err", xlog.String("errKind", "listen err"), xlog.FieldErr(err))
	}
	return &Server{
		Server:   newServer,
		listener: listener,
		config:   config,
	}
}

func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	return err
}

func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

func (s *Server) GracefulStop(_ context.Context) error {
	s.Server.GracefulStop()
	return nil
}

func (s *Server) Info() *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:    ModName,
		Scheme:  "grpc",
		Address: s.config.Addr,
	}
}
