package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tenkiej/config"
	"tenkiej/contoroller"
	"tenkiej/webserver"
	"time"

	"github.com/comail/colog"
)

func logConfig() error {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	return nil
}

func init() {
	if err := config.Init(); err != nil {
		log.Panic(err)
	}
	if err := logConfig(); err != nil {
		log.Panic(err)
	}
	if err := webserver.Init(); err != nil {
		log.Panic(err)
	}
	contoroller.Init()
}

func Stop(ctx context.Context) {
	contoroller.Stop(ctx)
	webserver.Stop(ctx)
}

func main() {
	ctxt := context.Background()
	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctxt)
	defer config.TracerStop(ctxt)
	ctx := context.Background()
	ctxs := context.Background()
	//シグナル
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go contoroller.Run(ctx)
	time.Sleep(time.Millisecond * 500)
	if err := contoroller.Wait(ctx); err != nil {
		log.Panic(err)
	}
	go webserver.Start(ctxs)
	<-sigs
	Stop(context.Background())
}
