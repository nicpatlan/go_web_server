package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (aCfg *ApiConfig) DeletePostHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get("Authorization")
	_, token, ok := strings.Cut(tokenStr, "Bearer ")
	if !ok {
		respondWithError(wr, http.StatusBadRequest, "invalid token header")
		return
	}
	authorID, err := ParseToken(token, aCfg.jwtsecret)
	if err != nil {
		respondWithError(wr, http.StatusUnauthorized, err.Error())
		return
	}
	postID, err := strconv.Atoi(req.PathValue("postID"))
	if err != nil {
		respondWithError(wr, http.StatusBadRequest, "Error converting postID")
		return
	}
	err = aCfg.database.DeletePost(postID, authorID)
	if err != nil {
		if err.Error() == "invalid authorID" {
			respondWithError(wr, http.StatusForbidden, err.Error())
			return
		}
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusNoContent, "")
}
