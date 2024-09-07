package govern

import (
	"esfgit.leju.com/golang/frame/server"
	"esfgit.leju.com/golang/frame/xlog"
	"golang.org/x/net/context"
	"net"
	"net/http"
)

type Server struct {
	*http.Server
	listener net.Listener
	config   *Config
}

func newServer(config *Config) *Server {
	var listener, err = net.Listen("tcp4", config.Addr)
	if err != nil {
		xlog.Error("govern start error", xlog.FieldErr(err))
	}

	return &Server{
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: DefaultServeMux,
		},
		listener: listener,
		config:   config,
	}
}

func (s *Server) Serve() error {
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
