package xmux

import (
	"context"
	"github.com/Remember9/frame/server"
	"github.com/Remember9/frame/xlog"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Server struct {
	*mux.Router
	Server   *http.Server
	config   *Config
	listener net.Listener
}

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		xlog.Error("new "+ModName+" server err", xlog.FieldErr(err))
	}
	return &Server{
		Router:   mux.NewRouter(),
		config:   config,
		listener: listener,
	}
}

func (s *Server) Serve() error {
	s.Server = &http.Server{
		Handler:      s.Router,
		Addr:         s.config.Addr,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Stop() error {
	return s.Server.Close()
}

func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) Info() *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:    ModName,
		Scheme:  "http",
		Address: s.config.Addr,
	}
}
