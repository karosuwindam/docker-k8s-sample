package main

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		defer wg.Done()
		if err := loop.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)
	if err := loop.RunWait(); err != nil {
		fmt.Println("Runloop wait timeout :", err)
	}
	go func() {
		defer wg.Done()
		if err := webserver.Start(); err != nil {
			panic(err)
		}
	}()

	<-sigs
	cancel()
	Stop()
	wg.Wait()
	return nil
}

func Stop() {
	loop.Stop()
	webserver.Stop()
	fmt.Println("main Shutdown")
}

func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	if err := Start(); err != nil {
		panic(err)
	}
}
