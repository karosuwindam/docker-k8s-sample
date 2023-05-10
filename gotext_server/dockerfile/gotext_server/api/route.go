package api

import (
	"gocsvserver/api/text"
	"gocsvserver/config"
	"gocsvserver/webserver"
)

var Route []webserver.WebConfig = []webserver.WebConfig{}

func Setup(cfg *config.Config) error {
	if r, err := text.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, r...)
	}
	return nil
}
