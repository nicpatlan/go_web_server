package main

import (
	"net/http"
)

type healthzHandler struct{}

func (healthzHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("OK"))
}
