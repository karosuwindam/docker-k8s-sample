package main

import (
	"context"
	"eijent/config"
	"eijent/controller"
	"eijent/webserver"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/comail/colog"
)

func Init() error {
	if err := config.Init(); err != nil {
		return err
	}
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: config.Log.Colors,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	if err := webserver.Init(); err != nil {
		return err
	}
	if err := controller.Init(); err != nil {
		return err
	}
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) { //センサー監視
		defer wg.Done()
		if err := controller.Run(ctx); err != nil {
			log.Panicln(err)
		}
	}(ctx)
	controller.Wait()

	go func(ctx context.Context) { //Webserver起動
		defer wg.Done()
		if err := webserver.Start(); err != nil {
			log.Panicln(err)
		}
	}(ctx)
	<-sigs
	Stop(ctx)
	wg.Wait()
	cancel()
	return nil
}

func Stop(ctx context.Context) {
	if err := webserver.Stop(ctx); err != nil {
		log.Println("errror:", err)
	}
	if err := controller.Stop(ctx); err != nil {
		log.Println("errror:", err)

	}
}

func main() {
	if err := Init(); err != nil {
		log.Panic("error:", err)
	}
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Panic("error:", err)
	}
}
