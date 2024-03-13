package main

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
)

func Init() error {
	if err := config.Init(); err != nil {
		return err
	}
	if err := loop.Init(); err != nil {
		return err
	}
	if err := webserver.Init(); err != nil {
		return err
	}
	return nil
}

func Start() error {
	webserver.Start()
	return nil
}

func Stop() {

}

func main() {
	if err := Init(); err != nil {
		panic(err)
	}
	if err := Start(); err != nil {
		panic(err)
	}
	Stop()
}
