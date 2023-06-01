package main

import (
	"context"
	"errors"
	"fmt"
	"ingresslist/api"
	"ingresslist/config"
	"ingresslist/getkube"
	"ingresslist/textroot"
	"ingresslist/webserver"
	"log"
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
	// kubeconfigの設定
	getkube.Setup(cfg)
	wcfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	//apiの設定 /v1/*
	if err := api.Setup(cfg); err != nil {
		return nil, err
	} else {
		webserver.Config(wcfg, api.Route, "/v1")
	}

	// textrootの設定
	if r, err := textroot.Setup(cfg); err != nil {
		return nil, err
	} else {
		webserver.Config(wcfg, r, "")
	}
	return wcfg.NewServer()
}

func Run() error {

	errCh := make(chan error, 1)
	errCh1 := make(chan error, 1)
	ctx := context.Background()
	//ストップシグナルを受け取るコンテキストを作成
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go app.srv.Run(ctx, errCh)
	go getkube.Run(ctx, errCh1)
	// シグナルを受け取ったら終了
	<-ctx.Done()
	if err := Shutdown(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	stop()
	select {
	case err := <-errCh:
		if len(errCh1) > 0 {
			fmt.Println("error")
			return err
		}
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
