package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondMarshallError(wr http.ResponseWriter, err string) {
	log.Printf("Error marshalling JSON: %s", err)
	wr.WriteHeader(http.StatusInternalServerError)
}

func respondWithError(wr http.ResponseWriter, statusCode int, msg string) {
	type invalidPost struct {
		Error string `json:"error"`
	}
	resBody := invalidPost{
		Error: msg,
	}
	invalidRes, err := json.Marshal(resBody)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}
	wr.WriteHeader(statusCode)
	wr.Write(invalidRes)
}

func respondWithJSON(wr http.ResponseWriter, statusCode int, payload interface{}) {
	type validPost struct {
		Body interface{} `json:"cleaned_body"`
	}
	resBody := validPost{
		Body: payload,
	}
	validRes, err := json.Marshal(resBody)
	if err != nil {
		respondMarshallError(wr, err.Error())
		return
	}
	wr.WriteHeader(statusCode)
	wr.Write(validRes)
}
