package main

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"context"
	"errors"
	"log/slog"
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		slog.InfoContext(ctx, "Server is shutting down...")
		defer cancel()
		Stop(ctx)
		slog.InfoContext(ctx, "Server is shut down")
		close(idleConnsClosed)
	}()
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tshutdown, terr := config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	if terr != nil {
		defer tshutdown(context.Background())
	}
	go func(ctx context.Context) {
		defer wg.Done()
		if err := loop.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)
	if err := loop.RunWait(); err != nil {
		slog.ErrorContext(ctx, "Runloop wait timeout :", err)
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
		slog.ErrorContext(ctx, "loop.Stop", err)
	}
	if err := webserver.Stop(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.ErrorContext(ctx, "HTTP server Shutdown: timeout")

		} else {
			slog.ErrorContext(ctx, "webserver.Stop", err)

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
	slog.Info("All Shutdown")
}
