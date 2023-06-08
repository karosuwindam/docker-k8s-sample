package main

import (
	"app/config"
	"app/webserver"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type APP struct {
	// Webサーバの管理関数
	srv *webserver.Server
}

var app *APP

func Config() (*webserver.Server, error) {
	cfg, err := config.Setup()
	if err != nil {
		return nil, err
	}
	wcfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	return wcfg.NewServer()
}

func Run() error {

	errCh := make(chan error, 1)
	ctx := context.Background()
	//ストップシグナルを受け取るコンテキストを作成
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go app.srv.Run(ctx, errCh)
	// シグナルを受け取ったら終了
	<-ctx.Done()
	Shutdown()
	stop()
	select {
	case err := <-errCh:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("timeout")
	}
}

func Shutdown() error {
	var err error = nil
	fmt.Println("shutdown")
	if e := app.srv.Wait(); e != nil {
		err = e
	}
	return err
}

func main() {
	srv, err := Config()
	if err != nil {
		panic(err)
	}
	app = &APP{
		srv: srv,
	}
	if err := Run(); err != nil {
		panic(err)
	}
}
