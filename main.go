package main

import (
	"fmt"
	"log"
	"net/http"
)

type healthzHandler struct{}

func (healthzHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
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

func (f *fileHits) GetHitsHandler() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Content-Type", "text/html; charset=utf-8")
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte("<html>\n<body>\n\t<h1>Welcome, Admin</h1>\n"))
		wr.Write([]byte(fmt.Sprintf("\t<p>Server has been visited %d times!</p>\n", f.fileserverHits)))
		wr.Write([]byte("</body>\n</html>"))
	})
}

func (f *fileHits) GetResetHandler() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		f.fileserverHits = 0
		wr.WriteHeader(http.StatusOK)
	})
}

func main() {
	// constants
	const filePattern = "/app/*"
	const fileStrip = "/app"
	const fileRoot = "."
	const healthzPattern = "GET /api/healthz"
	const metricPattern = "GET /admin/metrics"
	const resetPattern = "/api/reset"
	const port = "8080"

	// create server mux handler and fileHits counter
	serveMux := http.NewServeMux()
	fHits := fileHits{}

	// add handlers
	serveMux.Handle(filePattern, fHits.incrFileHits(http.StripPrefix(fileStrip, http.FileServer(http.Dir(fileRoot)))))
	serveMux.Handle(healthzPattern, healthzHandler{})
	serveMux.Handle(metricPattern, fHits.GetHitsHandler())
	serveMux.Handle(resetPattern, fHits.GetResetHandler())

	// create server on localhost port 8080
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// begin listening and responding to requests
	log.Printf("Running server from %s on port: %s", fileRoot, port)
	log.Fatal(server.ListenAndServe())
}
