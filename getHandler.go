package main

import (
	"net/http"

	"github.com/nicpatlan/go_web_server/internal/database"
)

type GetHandler struct {
	Database *database.DB
}

func (gh GetHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	posts, err := gh.Database.GetPosts()
	if err != nil {
		respondWithJSON(wr, http.StatusOK, database.Post{})
		return
	}
	respondWithJSON(wr, http.StatusOK, posts)
}
