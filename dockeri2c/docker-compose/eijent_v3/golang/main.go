package main

import (
	"context"
	"eijent/config"
	"eijent/controller"
	"eijent/webserver"
	"log"

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
	return nil
}

func main() {
	if err := Init(); err != nil {
		log.Panic("error:", err)
	}
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Panic("error:", err)
	}
	<-ctx.Done()
}
