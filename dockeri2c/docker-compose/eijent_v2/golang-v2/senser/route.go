package senser

import (
	"app/webserver"
	"fmt"
	"net/http"
)

var Route []webserver.WebConfig = []webserver.WebConfig{
	// {Pass: "/metrics", Handler: metrics},
	// {Pass: "/health", Handler: health},
	// {Pass: "/json", Handler: json},
	{Pass: "/", Handler: rootdate},
}

func rootdate(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}
