package main

import (
	"fmt"
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

type fileHits struct {
	fileserverHits int
}

func (f *fileHits) incrFileHits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		f.fileserverHits++
		handler.ServeHTTP(wr, req)
	})
}

func (f *fileHits) GetHitsHandler(path string) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(wr, req)
			return
		}
		wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte(fmt.Sprintf("Hits: %d", f.fileserverHits)))
	})
}

func (f *fileHits) GetResetHandler(path string) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(wr, req)
			return
		}
		f.fileserverHits = 0
		wr.WriteHeader(http.StatusOK)
	})
}

func main() {
	// constants
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "/healthz"
	const metricPattern = "/metrics"
	const resetPattern = "/reset"
	const port = "8080"

	// create server mux handler and fileHits counter
	serveMux := http.NewServeMux()
	fHits := fileHits{}

	// add handlers
	serveMux.Handle(filePattern, fHits.incrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, fHits.GetHitsHandler(metricPattern))
	serveMux.Handle(resetPattern, fHits.GetResetHandler(resetPattern))

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
