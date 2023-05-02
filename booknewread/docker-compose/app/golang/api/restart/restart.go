package restart

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"fmt"
	"net/http"
)

var apiname string = "restart"

func restart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + ":" + r.URL.Path)

	if r.Method == "POST" {
		loop.Reset_ON(loop.RESET_DATA)
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "NG")

	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: restart},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
