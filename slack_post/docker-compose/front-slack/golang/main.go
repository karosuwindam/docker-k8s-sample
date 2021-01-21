package main

import (
	"fmt"
	"os"
)

const (
	PORT = "8080"
)

func EnvSetup(server *Server) {
	if str := os.Getenv("WEB_PORT"); str != "" {
		server.Port = str
	}
	if str := os.Getenv("WEB_IP"); str != "" {
		server.Ip = str
	}
	if str := os.Getenv("SEND_URL"); str != "" {
		server.Url = str
	}
	if str := os.Getenv("SEND_PORT"); str != "" && server.Url != "" {
		server.Url += ":" + str
	}
}

func main() {
	server := Server{Port: PORT}
	EnvSetup(&server)
	if server.Url == "" {
		fmt.Println("Input not Send URL data set env SEND_URL")
		return
	}
	server.Start()
}
