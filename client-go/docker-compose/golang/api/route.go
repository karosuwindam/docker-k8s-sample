package api

import (
	"ingresslist/api/jsonget"
	"ingresslist/config"
	"ingresslist/webserver"
)

var Route []webserver.WebConfig = []webserver.WebConfig{}

func Setup(cfg *config.Config) error {
	if r, err := jsonget.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, r...)
	}
	return nil
}
