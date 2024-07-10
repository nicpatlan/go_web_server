package main

import (
	"encoding/json"
	"net/http"
)

type PostHandler struct{}

func (PostHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
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
	respondWithJSON(wr, http.StatusOK, newPost.Body)
}
