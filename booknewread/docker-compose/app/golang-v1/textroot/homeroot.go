package textroot

import "book-newread/webserver"

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/", Handler: viewhtml},
}
