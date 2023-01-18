package main

import (
	"app/webserver"
)

type Htmldata struct{}

func setupRoute() (Htmldata, error) {
	output := Htmldata{}
	return output, nil
}
func (t *Htmldata) setupRoute(cfg *webserver.SetupServer) {
}
