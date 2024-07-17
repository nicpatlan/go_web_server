package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (aCfg *ApiConfig) GetOnePostHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	postID, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		respondWithError(wr, http.StatusBadRequest, "Error converting postID")
		return
	}
	posts, err := aCfg.database.GetPosts()
	wr.Header().Set("Content-Type", "application/json")
	if err != nil || postID < 1 || postID > len(posts) {
		respondWithError(wr, http.StatusNotFound, fmt.Sprintf("No post with ID: %d", postID))
		return
	}
	respondWithJSON(wr, http.StatusOK, posts[postID-1])
}
