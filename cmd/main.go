package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/notifire/internal/server"
)

func main() {
	serverConfig := server.LoadConfig()
	server := server.New(serverConfig)
	if err := server.Start(); err != nil {
		log.Fatalf("server failed to start on port %d: %v", serverConfig.Port, err)
	}
}
