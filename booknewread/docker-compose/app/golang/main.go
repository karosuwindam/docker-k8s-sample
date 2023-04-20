package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Config() {

}

func Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	<-ctx.Done()
	cancel()
	return nil
}

func Shutdown() {
	fmt.Println("Shutdown")
}

func main() {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	Config()
	go Run(ctx)
	<-ctx.Done()
	stop()
	Shutdown()
}
