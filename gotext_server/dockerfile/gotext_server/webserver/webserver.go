package webserver

import (
	"context"
	"errors"
	"fmt"
	"gocsvserver/config"
	"log"
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

// Status
type Status struct {
	Status string `json:status`
}

var passList map[string]bool //登録したパスリスト

var sampleword string = "hello world" //サンプルテキスト

// hello(w http.ResponseWriter, r *http.Request)
//
// # サンプル用の取得コード　hello worldを返す
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, sampleword)
}

// NewSetup(*config.Config) = *SetupServer,error
//
// # Webサーバ設定の初期化関数
//
// data(*config.Config) : Env設定で読みだした設定
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

// (*SetupServer) NewServer() = *Server,error
//
// Webサーバの開始設定
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

// (*SetupServer) Add(route, handler) = error
//
// http上に定義されたハンドラ関数を紐づける 二重登録でエラーを返す
//
// route(string) : ホームからのURLルートパス
// handler(func(http.ResponseWriter, *http.Request)) : httpの関数処理
func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) error {
	if passList[route] {
		return errors.New("Added Root data :" + route)
	}
	passList[route] = true
	t.mux.HandleFunc(route, handler)
	return nil
}

// (*SetupServer) muxHandler()
// SetupServer内のmuxhandlerを返す関数
func (t *SetupServer) muxHandler() http.Handler { return t.mux }

// (s *Server) Run(ctx context.Context) = error
// サーバをスタートする関数
func (s *Server) Run(ctx context.Context, ch chan<- error) {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		fmt.Println("Start Server")
		if err := s.srv.Serve((s.l)); err != nil &&
			err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	<-ctx.Done()
	cancel()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	fmt.Println("shutdown")
	ch <- eg.Wait()
	return
}
