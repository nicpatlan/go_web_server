package main

import (
	"encoding/json"
	"net/http"
)

type TokenResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (aCfg ApiConfig) LoginUserHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	userLogin := UserRequest{}
	err := decoder.Decode(&userLogin)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	resp, err := aCfg.database.ValidateUserPassword(userLogin.Email, userLogin.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			respondWithError(wr, http.StatusUnauthorized, err.Error())
			return
		}
		if err.Error() == "user not found" {
			respondWithError(wr, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	resp.Token, err = GenerateUserToken(resp.ID, aCfg.jwtsecret)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, resp)
}
