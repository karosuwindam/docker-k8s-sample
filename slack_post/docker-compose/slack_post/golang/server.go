package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	Port  string
	Ip    string
	slack SlackApi
	api   ApiData
}

type message struct {
	Url     string `json:url`
	Message string `json:message`
}

func (t *Server) apiselect(w http.ResponseWriter, r *http.Request) {
	var tmp message
	tmp.Url = r.RequestURI
	t.api.InitSlack(t.slack)
	t.api.ApiSelect(w, r)
	tmp.Message = t.api.Message
	output, _ := json.Marshal(tmp)
	fmt.Fprintf(w, string(output)+"\n")
}

func (t *Server) hello(w http.ResponseWriter, r *http.Request) {
	output := "Hello World"
	fmt.Fprintf(w, output)
}

func (t *Server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.HandlerFunc(t.apiselect))
	mux.Handle("/help", http.HandlerFunc(t.hello))

	s := http.Server{
		Addr:    t.Ip + ":" + t.Port,
		Handler: mux,
	}

	fmt.Println(t.Ip + ":" + t.Port + " server start")
	s.ListenAndServe()
}
