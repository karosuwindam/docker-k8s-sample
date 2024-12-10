package main

import (
	"context"
	"gocsvserver/config"
	"gocsvserver/webserver"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Init() error {
	config.Init()
	webserver.Init()
	return nil
}

func Start() error {
	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	tshutdown, terr := config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	if terr != nil {
		defer tshutdown(context.Background())
	}
	ctxweb, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := webserver.Start(ctx); err != nil {
			panic(err)
		}
	}(ctxweb)

	<-sigs
	cancel()
	Stop()
	wg.Wait()
	return nil
}

func Stop() {
	webserver.Stop()
	slog.Info("main Shutdown")
}
func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	if err := Start(); err != nil {
		panic(err)
	}
}
