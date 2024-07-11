package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nicpatlan/go_web_server/internal/database"
)

type PostHandler struct {
	Database *database.DB
}

func (ph PostHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	type post struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	newPost := post{}
	err := decoder.Decode(&newPost)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	if len(newPost.Body) > 140 {
		respondWithError(wr, http.StatusBadRequest, "140 character limit exceeded")
		return
	}
	newPost.Body = cleanPost(newPost.Body)
	p, err := ph.Database.CreatePost(newPost.Body)
	if err != nil {
		log.Printf("Error creating post: %s", err)
		return
	}
	respondWithJSON(wr, http.StatusCreated, p)
}
