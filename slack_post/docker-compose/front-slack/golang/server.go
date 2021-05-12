package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	Port string
	Ip   string
	Url  string
}

type HealthMessage struct {
	Message string `json:message`
}

func (t *Server) postdata(w http.ResponseWriter, r *http.Request) {
	output := "Hello World2"
	fmt.Fprintf(w, output)
}

func (t *Server) hello(w http.ResponseWriter, r *http.Request) {
	output := "Hello World"
	fmt.Fprintf(w, output)
}

func (t *Server) health(w http.ResponseWriter, r *http.Request) {
	code := 200
	tmp := HealthMessage{Message: "OK"}

	url := "http://" + t.Url + "/health"

	resp, err := http.Get(url)
	if err != nil {
		// code = 204
		tmp.Message = "busy now"
		output, _ := json.Marshal(tmp)
		w.WriteHeader(code)
		fmt.Fprintf(w, string(output))
		fmt.Println(url + ":Time out")
		return
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(url, string(byteArray))

	w.WriteHeader(code)
	output, _ := json.Marshal(tmp)
	fmt.Fprintf(w, string(output))
}

func (t *Server) Start() {
	mux := http.NewServeMux()
	url, err := url.Parse("http://" + t.Url + "/api/v1")
	if err != nil {
		log.Println("Reverse Proxy target url could not be parsed:", err)
		return
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	mux.HandleFunc("/health", t.health)
	mux.Handle("/postmessage", httputil.NewSingleHostReverseProxy(url))

	s := http.Server{
		Addr:    t.Ip + ":" + t.Port,
		Handler: mux,
	}

	fmt.Println(t.Ip + ":" + t.Port + " server start")
	s.ListenAndServe()
}
