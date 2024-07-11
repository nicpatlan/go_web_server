package main

import (
	"encoding/json"
	"net/http"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (aCfg *ApiConfig) PostHandlerFunc(wr http.ResponseWriter, req *http.Request) {
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
	p, err := aCfg.database.CreatePost(newPost.Body)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, "Error creating post")
		return
	}
	respondWithJSON(wr, http.StatusCreated, p)
}
