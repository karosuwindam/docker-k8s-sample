package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//データのリロード
func (t *WebSetupData) reload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	} else {
		fmt.Fprintf(w, "Not Reload")
	}

}

//ステータスを取得する
func (t *WebSetupData) status(w http.ResponseWriter, r *http.Request) {
	jsondata, err := json.Marshal(GrobalStatus)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}
}
