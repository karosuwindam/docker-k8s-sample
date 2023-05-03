package main

import (
	"context"
	"fmt"
	"gocsvserver/config"
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
	wcfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	//ハンドラー登録設定
	// webserver.Config(wcfg, ,"")

	return wcfg.NewServer()
}

func Run() error {

	errCh := make(chan error, 1)
	ctx := context.Background()
	//ストップシグナルを受け取るコンテキストを作成
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	if s, err := Config(); err != nil {
		return err
	} else {
		go s.Run(ctx, errCh)
	}
	ctx.Done()
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
	if err := Run(); err != nil {
		log.Panicln(err)
		os.Exit(1)

	}

	fmt.Println("end")
}
