package main

import (
	"encoding/json"
	"net/http"
)

type UserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (aCfg *ApiConfig) CreateUserHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	newUser := UserRequest{}
	err := decoder.Decode(&newUser)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	resp, err := aCfg.database.CreateUser(newUser.Email, newUser.Password)
	if err != nil {
		if err.Error() == "duplicate email" {
			respondWithError(wr, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(wr, http.StatusInternalServerError, "Error creating user")
		return
	}
	respondWithJSON(wr, http.StatusCreated, resp)
}
