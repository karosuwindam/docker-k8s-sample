package jsonnobel

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"encoding/json"
	"fmt"
	"net/http"
)

var apiname string = "jsonnobel"

// ステータスを取得する

func getnowdata(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + ":" + r.URL.Path)

	if r.Method == "GET" {
		jsondata, err := json.Marshal(loop.ReadNListData())
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: getnowdata},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
