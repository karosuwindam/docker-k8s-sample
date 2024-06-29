package webserver

import (
	"book-newread/config"
	"book-newread/webserver/api"
	"book-newread/webserver/indexpage"
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

var srv *http.Server // Webサーバの管理関数

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
	api.Init(cfg.mux)
	// fileserver := http.FileServer(http.Dir(config.Web.StaticPage))
	// cfg.mux.Handle("/", fileserver)
	cfg.mux.HandleFunc("/", indexpage.Init("/"))
	return nil
}

func Start(ctx context.Context) error {
	var err error = nil
	if config.TraData.TracerUse {
		hander := otelhttp.NewHandler(cfg.mux, "http-server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		)
		srv = &http.Server{
			Addr:         cfg.hostname + ":" + cfg.port,
			Handler:      hander,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	} else {
		srv = &http.Server{
			Addr:         cfg.hostname + ":" + cfg.port,
			Handler:      cfg.mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	log.Println("info:", "Start Server", cfg.hostname+":"+cfg.port)
	go func() {
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
			panic(err)
		} else {
			err = nil
		}
	}()
	select {
	case <-shutdown:
		done <- true
		break
	case <-ctx.Done():
		return nil
	}
	return err
}

func Stop(ctx context.Context) error {
	if srv == nil {
		return nil
	}
	ctx, _ = context.WithTimeout(ctx, time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	shutdown <- true

	select {
	case <-done:
		break
	case <-ctx.Done():
		log.Println("ctx done")
		break
	case <-time.After(time.Microsecond * 500):
		log.Println("errors:", "shutdown time out over 500 ms")
		break
	}
	return nil
}
