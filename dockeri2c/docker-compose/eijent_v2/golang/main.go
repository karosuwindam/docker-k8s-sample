package main

import (
	"app/webserver"
	"context"
	"fmt"
	"log"
	"os"
)

func Config(cfg *webserver.SetupServer) (Htmldata, error) {
	h := setupRoute()
	h.setupRoute(cfg)
	return h, nil
}

func Run(ctx context.Context) error {
	cfg, err := webserver.NewSetup()
	if err != nil {
		return err
	}
	Config(cfg)

	return nil
}

func main() {
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
}
