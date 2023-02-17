package main

import (
	"fmt"
	"net/http"
)

type WebSetup struct {
	Ip       string `json:ip`
	Port     string `json:port`
	RootPath string `json:rootpath`
}

type WebSetupData struct {
	Data WebSetup
	Tmp  string
}

func (t *WebSetupData) metrics(w http.ResponseWriter, r *http.Request) {
	output := "senser_data{type=\"tmp\"} " + t.Tmp
	fmt.Fprintf(w, "%s", output)

}

func (t *WebSetupData) view(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)
}

func (t *WebSetupData) Webstart() {
	fmt.Println(t.Data.Ip + ":" + t.Data.Port + "server start")
	http.HandleFunc("/metrics", t.metrics)
	http.HandleFunc("/", t.view)
	http.ListenAndServe(t.Data.Ip+":"+t.Data.Port, nil)
}
