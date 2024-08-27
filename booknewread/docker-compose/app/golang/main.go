package main

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Init() error {
	if err := config.Init(); err != nil {
		return err
	}
	if err := loop.Init(); err != nil {
		return err
	}
	if err := webserver.Init(); err != nil {
		return err
	}
	return nil
}

func Start() error {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		//シャットダウン処理
		log.Println("info:", "Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		Stop(ctx)
		log.Println("info:", "Server is shut down")
		close(idleConnsClosed)
	}()
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	defer config.TracerStop(ctx)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := loop.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)
	if err := loop.RunWait(); err != nil {
		log.Println("error:", "Runloop wait timeout :", err)
	}
	if err := webserver.Start(ctx); err != nil {
		panic(err)
	}

	<-idleConnsClosed
	wg.Wait()
	return nil
}

func Stop(ctx context.Context) {
	if err := loop.Stop(context.Background()); err != nil {
		log.Println("error:", err)
	}
	if err := webserver.Stop(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("error:", "HTTP server Shutdown: timeout")

		} else {
			log.Println("error:", err)

		}
	}
}

func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	if err := Start(); err != nil {
		panic(err)
	}
	log.Println("info:", "All Shutdown")
}
