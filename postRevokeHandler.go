package main

import (
	"net/http"
	"strings"
)

func (aCfg *ApiConfig) RevokeTokenHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	refreshToken := req.Header.Get("Authorization")
	_, token, ok := strings.Cut(refreshToken, "Bearer ")
	if ok {
		err := aCfg.database.RevokeUserRefreshToken(token)
		if err != nil {
			if err.Error() == "invalid refresh token" {
				respondWithError(wr, http.StatusBadRequest, err.Error())
				return
			}
			respondWithError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		wr.WriteHeader(http.StatusNoContent)
		return
	}
	respondWithError(wr, http.StatusBadRequest, "invalid header token")
}
