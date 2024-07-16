package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Post struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (aCfg *ApiConfig) PostHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	type post struct {
		Body string `json:"body"`
	}

	tokenStr := req.Header.Get("Authorization")
	_, token, ok := strings.Cut(tokenStr, "Bearer ")
	wr.Header().Set("Content-Type", "application/json")
	if !ok {
		respondWithError(wr, http.StatusBadRequest, "invalid token header")
		return
	}
	userID, err := ParseToken(token, aCfg.jwtsecret)
	if err != nil {
		respondWithError(wr, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(req.Body)
	newPost := post{}
	err = decoder.Decode(&newPost)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}

	if len(newPost.Body) > 140 {
		respondWithError(wr, http.StatusBadRequest, "140 character limit exceeded")
		return
	}
	newPost.Body = cleanPost(newPost.Body)
	p, err := aCfg.database.CreatePost(userID, newPost.Body)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, "Error creating post")
		return
	}
	respondWithJSON(wr, http.StatusCreated, p)
}
