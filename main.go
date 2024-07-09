package main

import (
	"net/http"
)

func main() {
	// create server mux and handler
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	// create server on localhost port 8080
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// begin listening and responding to requests
	server.ListenAndServe()
}
