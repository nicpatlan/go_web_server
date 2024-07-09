package main

import (
	"log"
	"net/http"
)

type healthzHandler struct{}

func (healthzHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/healthz" {
		http.NotFound(wr, req)
		return
	}
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("OK"))
}

func main() {
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "/healthz"
	const port = "8080"

	// create server mux and handler
	serveMux := http.NewServeMux()
	serveMux.Handle(filePattern, http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot))))
	serveMux.Handle(healthzPattern, healthzHandler{})

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
