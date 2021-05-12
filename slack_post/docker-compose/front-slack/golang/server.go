package main

import (
	"fmt"
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

func (t *Server) postdata(w http.ResponseWriter, r *http.Request) {
	output := "Hello World2"
	fmt.Fprintf(w, output)
}

func (t *Server) hello(w http.ResponseWriter, r *http.Request) {
	output := "Hello World"
	fmt.Fprintf(w, output)
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
	mux.Handle("/postmessage", httputil.NewSingleHostReverseProxy(url))

	s := http.Server{
		Addr:    t.Ip + ":" + t.Port,
		Handler: mux,
	}

	fmt.Println(t.Ip + ":" + t.Port + " server start")
	s.ListenAndServe()
}
