package main

import (
	"flag"
	"log"

	"github.com/jekyll/dashboard"
)

func main() {
	var bindAddr string
	flag.StringVar(&bindAddr, "http", "localhost:8000", "The address (host:port) the server should listen on.")
	flag.Parse()

	log.Printf("Starting server on %s...", bindAddr)
	if err := dashboard.Listen(bindAddr); err != nil {
		log.Fatalf("Encountered error serving: %+v", err)
	}
}
