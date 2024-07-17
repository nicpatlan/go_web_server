package main

import (
	"net/http"
	"strconv"
)

func (aCfg *ApiConfig) GetPostsHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	authorIDStr := req.URL.Query().Get("author_id")
	sortStr := req.URL.Query().Get("sort")
	asc := false
	if sortStr == "" || sortStr == "asc" {
		asc = true
	}
	if authorIDStr == "" {
		posts, err := aCfg.database.GetPosts(0, false, asc)
		if err != nil {
			respondWithError(wr, http.StatusNotFound, "No posts to retrieve")
			return
		}
		respondWithJSON(wr, http.StatusOK, posts)
		return
	}
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	posts, err := aCfg.database.GetPosts(authorID, true, asc)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, posts)
}
