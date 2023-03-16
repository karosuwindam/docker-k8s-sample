package main

import (
	"app/senser"
	"app/webserver"
	"context"
	"fmt"
	"log"
	"os"
)

func RootConfg() []webserver.WebConfig {
	output := []webserver.WebConfig{}
	for _, route := range senser.Route {
		output = append(output, route)
	}
	return output
}

func Config(cfg *webserver.SetupServer) error {
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
	s, err := cfg.NewServer()
	if err != nil {
		return err
	}

	return s.Run(ctx)
}
func EndRun() {}
func main() {
	fmt.Println("start")
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
}
