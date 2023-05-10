package textroot

import "gocsvserver/webserver"

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/", Handler: viewhtml},
}
