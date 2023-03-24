package main

import (
	"flag"
	"log"

	"github.com/mkrtychanr/avito_backend_internship/internal/server"
)

func main() {
	flag.Parse()
	configPath := flag.Arg(0)
	server, err := server.NewServer(configPath)
	if err != nil {
		log.Fatal(err)
	}
	server.Up()
}
