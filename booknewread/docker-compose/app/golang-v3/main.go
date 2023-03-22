package main

import (
	"booknewread/loop"
	"booknewread/webserver"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	url := "http://ncode.syosetu.com/n3289ds/?p=4"
	url1 := "https://ncode.syosetu.com/n8920ex/?p=3"
	url2 := "https://circle.ms/"
	go func() {
		loop.Loop([]string{url, url1, url2})
	}()
	time.Sleep(time.Millisecond * 50 * 0)
	loop.Read()
	fmt.Println(loop.Count(), loop.ListData)

	return
	fmt.Println("start")
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	EndRun()
	fmt.Println("end")
}
