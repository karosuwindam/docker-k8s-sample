package main

import (
	"context"
	"errors"
	"ingresslist/config"
	"ingresslist/controller"
	"ingresslist/webserver"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/m-mizutani/clog"
)

func logconfig() {
	handler := clog.New(
		clog.WithColor(true),
		clog.WithSource(true),
	)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Init() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := controller.Init(); err != nil {
		panic(err)
	}
	if err := webserver.Init(); err != nil {
		panic(err)
	}
	logconfig()
}

func Stop() {
	if err := controller.Stop(context.Background()); err != nil {
		slog.Error("controller stop", "error", err)
	}
	if err := webserver.Stop(context.Background()); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("webserver Stop", "error", "HTTP server Shutdown: timeout")

		} else {
			slog.Error("webserver Stop", "error", err)

		}

	}

}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	Init()
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	defer config.TracerStop(ctx)

	go controller.Run(ctx)
	controller.Wait()
	go webserver.Start(ctx)
	<-sigs
	Stop()
}
