package main

import (
	"context"
	"fmt"
	"gocsvserver/api"
	"gocsvserver/config"
	"gocsvserver/textroot"
	"gocsvserver/webserver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Config() (*webserver.Server, error) {
	cfg, err := config.Setup()
	if err != nil {
		return nil, err
	}
	if err := api.Setup(cfg); err != nil {
		return nil, err
	}
	wcfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	//ハンドラー登録設定
	webserver.Config(wcfg, api.Route, "/v1")
	webserver.Config(wcfg, textroot.Route, "/")

	return wcfg.NewServer()
}

func Run(ctx context.Context) error {

	errCh := make(chan error, 1)
	//ストップシグナルを受け取るコンテキストを作成
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	if s, err := Config(); err != nil {
		return err
	} else {
		go s.Run(ctx, errCh)
	}
	<-ctx.Done()
	stop()
ErrCK:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-time.After(5 * time.Second):
			break ErrCK
		}
	}
	return Shutdown()
}

// シャットダウン中の処理
func Shutdown() error {
	return nil
}

func main() {
	// flag.Parse() //コマンドラインオプションの有効
	log.SetFlags(log.Llongfile | log.Flags())
	ctx := context.Background()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-sigChan
		fmt.Println("signal")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}()
	if err := Run(ctx); err != nil {
		log.Panicln(err)
		os.Exit(1)
	}

	fmt.Println("end")
}
