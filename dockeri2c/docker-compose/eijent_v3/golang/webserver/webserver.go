package webserver

import (
	"context"
	"eijent/config"
	"eijent/webserver/health"
	"eijent/webserver/jsonout"
	"eijent/webserver/metricsout"
	"eijent/webserver/reset"
	"eijent/webserver/rootpage"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

// SetupServer
// サーバ動作の設定
type SetupServer struct {
	protocol string // Webサーバーのプロトコル
	hostname string //Webサーバのホスト名
	port     string //Webサーバの解放ポート

	mux *http.ServeMux //webサーバのmux
}

var srv *http.Server // Webサーバの管理関数

var cfg SetupServer

func HelloWeb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Web"))
}

func APIInit(mux *http.ServeMux) error {
	mux.HandleFunc("GET /", HelloWeb)
	return nil
}

type api struct {
	Router string
	Func   func(string, *http.ServeMux) error
}

var routes = []api{
	{"/metrics", metricsout.Init},
	{"/json", jsonout.Init},
	{"/reset", reset.Init},
	{"/health", health.Init},
	{"/", rootpage.Init},
}

func Init() error {
	cfg = SetupServer{
		protocol: config.Web.Protocol,
		hostname: config.Web.Hostname,
		port:     config.Web.Port,
		mux:      http.NewServeMux(),
	}
	for _, r := range routes {
		if err := r.Func(r.Router, cfg.mux); err != nil {
			return errors.Wrapf(err, "setup %v", r.Router)
		}
	}
	return nil
}

func Start() error {
	var err error = nil
	var wg sync.WaitGroup
	srv = &http.Server{
		Addr:    cfg.hostname + ":" + cfg.port,
		Handler: cfg.mux,
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	log.Println("info: Start Server", cfg.hostname+":"+cfg.port)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
			panic(err)
		} else {
			err = nil
		}
	}()
	wg.Wait()
	log.Println("info: Server Stop")
	return err
}

func Stop(ctx context.Context) error {
	if srv == nil {
		return nil
	}
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
