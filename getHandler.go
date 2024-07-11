package main

import (
	"net/http"
)

func (aCfg *ApiConfig) GetPostsHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	posts, err := aCfg.database.GetPosts()
	if err != nil {
		respondWithError(wr, http.StatusNotFound, "No posts to retrieve")
		return
	}
	respondWithJSON(wr, http.StatusOK, posts)
}
