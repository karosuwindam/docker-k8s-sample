package webserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"golang.org/x/sync/errgroup"
)

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`

	mux *http.ServeMux
}

type Server struct {
	srv *http.Server
	l   net.Listener
}

type Status struct {
	Status string `json:status`
}

func NewSetup() (*SetupServer, error) {
	cfg := &SetupServer{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	cfg.mux = http.NewServeMux()
	return cfg, nil
}

func (t *SetupServer) NewServer() (*Server, error) {
	fmt.Println("Setupserver", t.Protocol, t.Hostname+":"+t.Port)
	l, err := net.Listen(t.Protocol, t.Hostname+":"+t.Port)
	if err != nil {
		return nil, err
	}
	return &Server{
		srv: &http.Server{Handler: t.muxHandler()},
		l:   l,
	}, nil
}

func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) {
	t.mux.HandleFunc(route, handler)
}

func (t *SetupServer) muxHandler() http.Handler { return t.mux }

func (s *Server) Run(ctx context.Context) error {
	// ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// defer stop()
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		fmt.Println("Start Server")
		if err := s.srv.Serve((s.l)); err != nil {
			return err
		}
		return nil
	})
	<-ctx.Done()
	cancel()
	return eg.Wait()
}

func (s *Server) Shutdown() error {
	err := s.srv.Shutdown(context.Background())
	fmt.Println("shutdown")
	return err

}
