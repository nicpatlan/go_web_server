package main

import (
	"net/http"
)

func (aCfg *ApiConfig) GetPostsHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	posts, err := aCfg.database.GetPosts()
	if err != nil {
		respondWithJSON(wr, http.StatusOK, Post{})
		return
	}
	respondWithJSON(wr, http.StatusOK, posts)
}
