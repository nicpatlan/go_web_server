package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (aCfg *ApiConfig) UpdateUserHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	bearerToken := req.Header.Get("Authorization")
	_, respToken, ok := strings.Cut(bearerToken, "Bearer ")
	if ok {
		id, err := ParseToken(respToken, aCfg.jwtsecret)
		if err != nil {
			respondWithError(wr, http.StatusUnauthorized, err.Error())
			return
		}
		decoder := json.NewDecoder(req.Body)
		userRequest := UserRequest{}
		decoder.Decode(&userRequest)
		userResponse, err := aCfg.database.UpdateUser(id, userRequest.Email, userRequest.Password)
		if err != nil {
			respondWithError(wr, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(wr, http.StatusOK, userResponse)
		return
	}
	respondWithError(wr, http.StatusBadRequest, "invalid header token")
}
