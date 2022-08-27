package webserver

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	log.Println("Setupserver", t.Protocol, t.Hostname+":"+t.Port)
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
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		log.Println("Start Server")
		if err := s.srv.Serve((s.l)); err != nil {
			return err
		}
		return nil
	})
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	log.Println("shutdown")
	return eg.Wait()
}
