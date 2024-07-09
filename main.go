package main

import (
	"net/http"
)

func main() {
	// create server handler
	serveMux := http.NewServeMux()

	// create server on localhost port 8080
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// begin listening and responding to requests
	server.ListenAndServe()
}
