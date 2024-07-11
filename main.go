package main

import (
	"log"
	"net/http"

	"github.com/nicpatlan/go_web_server/internal/database"
)

func main() {
	// constants
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "GET /api/healthz"
	const metricPattern = "GET /admin/metrics"
	const resetPattern = "/api/reset"
	const postPattern = "POST /api/posts"
	const getPattern = "GET /api/posts"
	const port = "8080"
	const dbPath = "database.json"

	// database setup
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("Error loading database: %s", err)
		return
	}

	// create server mux handler and fileHits counter
	serveMux := http.NewServeMux()
	fHits := FileHits{}

	// add handler
	serveMux.Handle(filePattern, fHits.IncrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, fHits.GetHitsHandler())
	serveMux.Handle(resetPattern, fHits.GetResetHandler())
	serveMux.Handle(postPattern, PostHandler{Database: db})
	serveMux.Handle(getPattern, GetHandler{Database: db})

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
