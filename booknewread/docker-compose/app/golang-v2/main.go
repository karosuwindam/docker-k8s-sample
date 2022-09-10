package main

import (
	"booknewread/webserver"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := webserver.NewSetup()
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}
	cc, err := setupbaseRoute()
	if err != nil {
		return
	}
	go func() {
		ckNobelloop(cc.Loopdata)
	}()
	go func() {
		ckBooklloop(cc.Loopdata)
	}()
	loopWait(cc.Loopdata)
	cc.setupRoute(cfg)
	s, err := cfg.NewServer()
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return

	}
	ctx := context.Background()
	if err := s.Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	loopStop(cc.Loopdata)
	fmt.Println("end")

}
