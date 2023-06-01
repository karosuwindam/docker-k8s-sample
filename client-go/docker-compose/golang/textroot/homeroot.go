package textroot

import (
	"ingresslist/config"
	"ingresslist/webserver"
)

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/", Handler: viewhtml},
}

func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return Route, nil
}
