package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UserHook struct {
	UserID int `json:"user_id"`
}

type Webhook struct {
	Event string   `json:"event"`
	Data  UserHook `json:"data"`
}

func (aCfg *ApiConfig) WebhookHandlerFunc(wr http.ResponseWriter, req *http.Request) {
	keyStr := req.Header.Get("Authorization")
	_, key, ok := strings.Cut(keyStr, "ApiKey ")
	if !ok || key != aCfg.apikey {
		respondWithError(wr, http.StatusUnauthorized, "invalid key header")
		return
	}
	decoder := json.NewDecoder(req.Body)
	webhook := Webhook{}
	err := decoder.Decode(&webhook)
	if err != nil {
		respondWithError(wr, http.StatusBadRequest, err.Error())
		return
	}
	if webhook.Event == "user.upgraded" {
		err = aCfg.database.UpdateUserPremium(webhook.Data.UserID, true)
		if err != nil {
			if err.Error() == "invalid userID" {
				respondWithError(wr, http.StatusNotFound, err.Error())
				return
			}
			respondWithError(wr, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondWithJSON(wr, http.StatusNoContent, "")
}
