package status

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"encoding/json"
	"fmt"
	"net/http"
)

var apiname string = "status"

// ステータスを取得する
func status(w http.ResponseWriter, r *http.Request) {
	jsondata, err := json.Marshal(loop.ReadStatus())
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: status},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
