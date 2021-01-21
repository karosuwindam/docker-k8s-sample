package main

import (
	"fmt"
	"os"
)

func EnvSetup(server *Server) {
	if str := os.Getenv("WEB_PORT"); str != "" {
		server.Port = str
	}
	if str := os.Getenv("WEB_IP"); str != "" {
		server.Ip = str
	}
	if str := os.Getenv("SLACK_TOKEN"); str != "" {
		server.slack.Token = str
	} else {
		server.slack.Token = TOKEN
	}
	if str := os.Getenv("SLACK_CHANNEL"); str != "" {
		server.slack.Channel = str
	} else {
		server.slack.Channel = CHANNEL
	}

}

func main() {
	server := Server{Port: "8080"}
	EnvSetup(&server)
	if server.slack.Token == "" {
		fmt.Println("Input not TOKEN data set env SLACK_TOKEN")
		return
	}
	server.Start()
}
