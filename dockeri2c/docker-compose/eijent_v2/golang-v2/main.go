package main

import (
	"app/senser"
	"app/uptime"
	"app/webserver"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

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
		fmt.Println("Start Read Sennser")
		for {
			senser.SenserRead()
			fmt.Println(senser.SennserDataValue)
			time.Sleep(5 * time.Second)
		}
		return nil
	})
	if senser.SennserData.Bme280_data.Flag {
		fmt.Println(senser.SennserData.Bme280_data.ReadData())
	}
	s, err := cfg.NewServer()
	if err != nil {
		return err
	}

	return s.Run(ctx)
}
func EndRun() {}

func main() {
	count := 0
	for {
		ck := uptime.Read()
		if ck > 120 {
			break
		} else {
			if count == 0 {
				fmt.Println("sleep wake up time", ck)
			}
			time.Sleep(time.Second)
		}
		count++
	}
	fmt.Println(uptime.Read())
	fmt.Println("start")
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
}
