package main

import (
	"fmt"
	"net/http"
)

type FileHits struct {
	fileserverHits int
}

func (f *FileHits) IncrFileHits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		f.fileserverHits++
		handler.ServeHTTP(wr, req)
	})
}

func (f *FileHits) GetHitsHandler() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Content-Type", "text/html; charset=utf-8")
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte("<html>\n<body>\n\t<h1>Welcome, Admin</h1>\n"))
		wr.Write([]byte(fmt.Sprintf("\t<p>Server has been visited %d times!</p>\n", f.fileserverHits)))
		wr.Write([]byte("</body>\n</html>"))
	})
}

func (f *FileHits) GetResetHandler() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		f.fileserverHits = 0
		wr.WriteHeader(http.StatusOK)
	})
}
