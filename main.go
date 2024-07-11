package main

import (
	"log"
	"net/http"

	"github.com/nicpatlan/go_web_server/internal/database"
)

type ApiConfig struct {
	fileserverHits int
	database       *database.DB
}

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
	const getOnePattern = "GET /api/posts/{id}"
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
	aCfg := ApiConfig{fileserverHits: 0, database: db}

	// add handler
	serveMux.Handle(filePattern, aCfg.IncrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, aCfg.GetHitsHandler())
	serveMux.Handle(resetPattern, aCfg.GetResetHandler())
	serveMux.HandleFunc(postPattern, aCfg.PostHandlerFunc)
	serveMux.HandleFunc(getPattern, aCfg.GetPostsHandlerFunc)
	serveMux.HandleFunc(getOnePattern, aCfg.GetOnePostHandlerFunc)

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
