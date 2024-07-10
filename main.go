package main

import (
	"log"
	"net/http"
)

func main() {
	// constants
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "GET /api/healthz"
	const metricPattern = "GET /admin/metrics"
	const resetPattern = "/api/reset"
	const postPattern = "POST /api/validate_post"
	const port = "8080"

	// create server mux handler and fileHits counter
	serveMux := http.NewServeMux()
	fHits := FileHits{}

	// add handlers
	serveMux.Handle(filePattern, fHits.IncrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, fHits.GetHitsHandler())
	serveMux.Handle(resetPattern, fHits.GetResetHandler())
	serveMux.Handle(postPattern, PostHandler{})

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
