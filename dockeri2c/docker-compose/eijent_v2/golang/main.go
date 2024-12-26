package main

import (
	"app/config"
	"app/senser"
	"app/uptime"
	"app/webserver"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func init() {

	if err := config.Init(); err != nil {
		panic(err)
	}
}

func RootConfg() []webserver.WebConfig {
	output := []webserver.WebConfig{}
	for _, route := range senser.Route {
		output = append(output, route)
	}
	return output
}

func Config(cfg *webserver.SetupServer) error {
	senser.SennserSetup()
	webserver.Config(cfg, RootConfg())
	return nil
}

func Run(ctx context.Context) error {
	cfg, err := webserver.NewSetup()
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	if err != nil {
		return err
	}
	if err := Config(cfg); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		slog.InfoContext(ctx,
			"Start Read Senser",
		)
		chdata := false
		for {
			senser.SenserRead()
			if !chdata {
				chdata = true
				ch1 <- true
			}
			time.Sleep(1 * time.Second)
		}
		return nil
	})
	eg.Go(func() error {
		slog.InfoContext(
			ctx,
			"Start Read Gyro",
		)
		chdata := false
		for {
			senser.SenserMoveRead()
			if !chdata {
				chdata = true
				ch2 <- true
			}
			time.Sleep((senser.GYRO_SLEEP_TIME))
		}
		return nil
	})

	s, err := cfg.NewServer()
	if err != nil {
		return err
	}
	<-ch1
	<-ch2
	return s.Run(ctx)
}
func EndRun() {
	senser.Close()
}

func main() {
	count := 0
	tshutdown, terr := config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	if terr != nil {
		defer tshutdown(context.Background())
	}
	for {
		ck := uptime.Read()
		if ck > 180 {
			break
		} else {
			if count == 0 {
				slog.Info(
					fmt.Sprintf("sleep wake up time %v", ck),
				)
			}
			time.Sleep(time.Second)
		}
		count++
	}
	slog.Info(
		fmt.Sprintf("uptime over %v starting", uptime.Read()),
	)
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	EndRun()
	slog.Info(
		fmt.Sprintf("senser stop end"),
	)
}
