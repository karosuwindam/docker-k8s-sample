package webserver

import (
	"context"
	"fmt"
	"gocsvserver/config"
	"gocsvserver/webserver/api"
	"gocsvserver/webserver/indexpage"
	"log/slog"
	"net"
	"net/http"
	"time"
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

var shutdown chan bool
var done chan bool

func HelloWeb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Web"))
}

func Init() error {
	shutdown = make(chan bool, 1)
	done = make(chan bool, 1)
	cfg = SetupServer{
		protocol: config.Web.Protocol,
		hostname: config.Web.Hostname,
		port:     config.Web.Port,
		mux:      http.NewServeMux(),
	}
	if err := api.Init(cfg.mux); err != nil {
		slog.Error("api.Init error", "error", err)
	}
	config.TraceHttpHandleFunc(cfg.mux, "/", indexpage.Init("/"))
	return nil
}

func Start(ctx context.Context) error {
	var err error = nil
	srv := &http.Server{
		Addr:         cfg.hostname + ":" + cfg.port,
		Handler:      cfg.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Start Server"+cfg.hostname+":"+cfg.port, "hostname", cfg.hostname, "port", cfg.port)
	go func() {
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {

		} else {
			err = nil
		}
	}()
	select {
	case <-shutdown:
		done <- true
		break
	case <-ctx.Done():
		return ctx.Err()
	}
	return err
}

func Stop() error {
	if len(shutdown) > 0 {
		return nil
	}
	shutdown <- true
	select {
	case <-done:
		slog.Debug("Web server stop")
		break
	case <-time.After(10 * time.Second):
		return fmt.Errorf("Web server stop timeout")
	}

	return nil
}
