package main

import (
	"encoding/json"
	"net/http"
)

func (aCfg *ApiConfig) CreateUserHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	type email struct {
		Body string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	newUser := email{}
	err := decoder.Decode(&newUser)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	p, err := aCfg.database.CreateUser(newUser.Body)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, "Error creating user")
		return
	}
	respondWithJSON(wr, http.StatusCreated, p)
}
