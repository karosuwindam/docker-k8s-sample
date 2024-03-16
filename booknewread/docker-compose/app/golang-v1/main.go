package main

import (
	"book-newread/api"
	"book-newread/config"
	"book-newread/loop"
	"book-newread/pyroscopesetup"
	"book-newread/textroot"
	"book-newread/webserver"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Config(cfg *config.Config) (*webserver.SetupServer, error) {
	api.Setup(cfg)
	loop.Setup(cfg)
	scfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	webserver.Config(scfg, api.Route, "/v1")
	webserver.Config(scfg, textroot.Route, "/")
	return scfg, nil
}

func Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	serverch := make(chan error)
	bookch := make(chan error)
	commonch := make(chan error)
	nobelch := make(chan error)
	var wgg sync.WaitGroup
	wgg.Add(1)
	go func() {
		defer wgg.Done()
		loop.BookloopStart()
	}()
	if cfg, err := config.EnvRead(); err != nil {
		cancel()
		return
	} else {
		if scfg, err := Config(cfg); err == nil {
			if s, err := scfg.NewServer(); err != nil {
				cancel()
				return
			} else {
				go loop.BookLoop(ctx, bookch)
				go loop.CommonLoop(ctx, commonch)
				go loop.NobelLoop(ctx, nobelch)
				wgg.Wait()
				go s.Run(ctx, serverch)
				<-ctx.Done()
				cancel()
				if err := s.Shutdown(); err != nil {
					log.Println(err)
				}
			}

		} else {
			cancel()
			return
		}
	}
	go func() {
		select {
		case <-time.After(time.Second):
			log.Println("time out")
			close(serverch)
			close(bookch)
			close(commonch)
			close(nobelch)
		}
	}()
	if err := <-serverch; err != nil {
		log.Println(err)
	}
	if err := <-bookch; err != nil {
		log.Println(err)
	}
	if err := <-commonch; err != nil {
		log.Println(err)
	}
	if err := <-nobelch; err != nil {
		log.Println(err)
	}
	return
}

func Shutdown() {
	fmt.Println("main Shutdown")
}

func main() {
	pyro := pyroscopesetup.Setup()
	pyroscopesetup.Add("base", "v1")
	pyro.Run()
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(1)
	go Run(ctx, &wg)
	<-ctx.Done()
	stop()
	wg.Done()
	Shutdown()
}
