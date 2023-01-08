package main

import (
	"log"
	"os"

	"github.com/dgparker/lilproxy/pkg/udpproxy"
)

func main() {
	target := os.Getenv("LILPROXY_TARGET")
	if target == "" {
		log.Fatal("env LILPROXY_TARGET required")
	}

	port := os.Getenv("LILPROXY_PORT")
	if port == "" {
		log.Fatal("env LILPROXY_PORT required")
	}

	c, err := udpproxy.New(port, target)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("lilproxy initialized")

	log.Fatal(c.ListenAndServe())
}
