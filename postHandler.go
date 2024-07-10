package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type PostHandler struct{}

func (PostHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	type post struct {
		Body string `json:"body"`
	}

	type validPost struct {
		Valid bool `json:"valid"`
	}

	type invalidPost struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	newPost := post{}
	err := decoder.Decode(&newPost)
	if err != nil {
		log.Printf("Error decoding post request: %s", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	if len(newPost.Body) > 140 {
		resBody := invalidPost{
			Error: "Something went wrong",
		}
		invalidRes, err := json.Marshal(resBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}
		wr.WriteHeader(http.StatusBadRequest)
		wr.Write(invalidRes)
		return
	}
	resBody := validPost{
		Valid: true,
	}
	validRes, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wr.WriteHeader(http.StatusOK)
	wr.Write(validRes)
}
