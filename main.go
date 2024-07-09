package main

import (
	"log"
	"net/http"
)

func main() {
	const fileRoot = "."
	const port = "8080"

	// create server mux and handler
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(fileRoot)))

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
