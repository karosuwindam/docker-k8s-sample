package webserver

import (
	"book-newread/config"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"golang.org/x/sync/errgroup"
)

// SetupServer
// サーバ動作の設定
type SetupServer struct {
	protocol string // Webサーバーのプロトコル
	hostname string //Webサーバのホスト名
	port     string //Webサーバの解放ポート

	mux *http.ServeMux //webサーバのmux
}

// Server
// Webサーバの管理情報
type Server struct {
	// Webサーバの管理関数
	srv *http.Server
	// 解放の管理関数
	l net.Listener
}

type Status struct {
	Status string `json:status`
}

var passList map[string]bool //登録したパスリスト

func NewSetup(data *config.Config) (*SetupServer, error) {
	cfg := &SetupServer{
		protocol: data.Server.Protocol,
		hostname: data.Server.Hostname,
		port:     data.Server.Port,
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	passList = map[string]bool{}
	cfg.mux = http.NewServeMux()
	return cfg, nil
}

func (t *SetupServer) NewServer() (*Server, error) {
	fmt.Println("Setup server", t.protocol, t.hostname+":"+t.port)
	l, err := net.Listen(t.protocol, t.hostname+":"+t.port)
	if err != nil {
		return nil, err
	}
	return &Server{
		srv: &http.Server{Handler: t.muxHandler()},
		l:   l,
	}, nil
}

func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) error {
	if passList[route] {
		return errors.New("Added Root data :" + route)
	}
	passList[route] = true
	t.mux.HandleFunc(route, handler)
	return nil
}

func (t *SetupServer) muxHandler() http.Handler { return t.mux }

func (s *Server) Run(ctx context.Context, ch chan<- error) {
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
	ch <- eg.Wait()
}

func (s *Server) Shutdown() error {
	err := s.srv.Shutdown(context.Background())
	fmt.Println("Server shutdown")
	return err

}
