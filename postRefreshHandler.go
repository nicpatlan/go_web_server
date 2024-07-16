package main

import (
	"net/http"
	"strings"
)

type SingleTokenReponse struct {
	Token string `json:"token"`
}

func (aCfg *ApiConfig) RefreshTokenHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	tokenResponse := SingleTokenReponse{}
	refreshToken := req.Header.Get("Authorization")
	_, token, ok := strings.Cut(refreshToken, "Bearer ")
	wr.Header().Set("Content-Type", "application/json")
	if ok {
		user, err := aCfg.database.ValidateUserRefreshToken(token)
		if err != nil {
			if err.Error() == "invalid refresh token" || err.Error() == "refresh token expired" {
				respondWithError(wr, http.StatusUnauthorized, err.Error())
				return
			}
			respondWithError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		token, err := GenerateUserToken(user.ID, aCfg.jwtsecret)
		if err != nil {
			respondWithError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		tokenResponse.Token = token
		respondWithJSON(wr, http.StatusOK, tokenResponse)
		return
	}
	respondWithError(wr, http.StatusBadRequest, "invalid header token")
}
