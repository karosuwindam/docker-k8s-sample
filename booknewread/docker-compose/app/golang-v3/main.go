package main

import (
	"context"
	"fmt"
	"gowebserver/webserver"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello World")
}

func RootConfg() []webserver.WebConfig {
	output := []webserver.WebConfig{}
	tmp := webserver.WebConfig{Pass: "/", Handler: hello}
	output = append(output, tmp)
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
