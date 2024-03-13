package webserver

import (
	"book-newread/config"
	"book-newread/webserver/api"
	"context"
	"net"
	"net/http"
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

var cfg SetupServer

var ctx context.Context
var cancel context.CancelFunc

func HelloWeb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Web"))
}

func Init() error {
	cfg = SetupServer{
		protocol: config.Web.Protocol,
		hostname: config.Web.Hostname,
		port:     config.Web.Port,
		mux:      http.NewServeMux(),
	}
	ctx, cancel = context.WithCancel(context.Background())
	api.Init(cfg.mux)
	return nil
}

func Start() error {
	var err error = nil
	srv := &http.Server{
		Addr:    cfg.hostname + ":" + cfg.port,
		Handler: cfg.mux,
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	go func() {
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	<-ctx.Done()
	return err
}

func Stop() error {
	cancel()
	return nil
}
