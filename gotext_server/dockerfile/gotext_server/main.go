package main

import (
	"context"
	"fmt"
	"gocsvserver/config"
	"gocsvserver/webserver"
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
	_, cancel := context.WithCancel(context.Background())
	wg.Add(1)
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
