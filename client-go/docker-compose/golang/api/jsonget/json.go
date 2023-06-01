package jsonget

import (
	"fmt"
	"ingresslist/config"
	"ingresslist/getkube"
	"ingresslist/webserver"
	"net/http"
)

var apiname = "json"

// httpハンドラ
func jsonget(w http.ResponseWriter, r *http.Request) {
	jdata := getkube.GetJsonData()
	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jdata)
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: jsonget},
}

func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
