package api

import (
	"book-newread/api/common"
	"book-newread/api/jsonb"
	"book-newread/api/jsonnobel"
	"book-newread/api/restart"
	"book-newread/api/status"
	"book-newread/config"
	"book-newread/webserver"
)

var Route []webserver.WebConfig = []webserver.WebConfig{}

func Setup(cfg *config.Config) error {

	//common
	if err := common.Setup(cfg); err != nil {
		return err
	}

	//status
	if tmp, err := status.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//restart
	if tmp, err := restart.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//jsonb
	if tmp, err := jsonb.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//jsonnobel
	if tmp, err := jsonnobel.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}

	return nil
}
