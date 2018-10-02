package main

import (
	"log"
	"net"
	"net/http"

	"github.com/tvpsh2020/messagebird-server/config"
	"github.com/tvpsh2020/messagebird-server/router"
	"github.com/tvpsh2020/messagebird-server/taskmanager"
)

func init() {
	config.Initialize()
	taskmanager.Initialize()
	router.Initialize()
}

func main() {
	listener, err := net.Listen("tcp", ":"+config.ServerPort)

	if err != nil {
		log.Printf("something went wrong when starting server: %s", err.Error())
		return
	}

	log.Printf("server start at port %s with %s mode", config.ServerPort, config.ServerMode)
	http.Serve(listener, nil)
}
