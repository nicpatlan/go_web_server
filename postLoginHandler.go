package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
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

	defaultExpiration := 60 * 60 * 24
	if userLogin.Expires == 0 {
		userLogin.Expires = defaultExpiration
	} else if userLogin.Expires > defaultExpiration {
		userLogin.Expires = defaultExpiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "go_web_app",
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(userLogin.Expires))),
		//ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(60*60*24) * time.Second)),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Subject:  strconv.Itoa(resp.ID),
	})
	tokenStr, err := token.SignedString([]byte(aCfg.jwtsecret))
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Token = tokenStr
	respondWithJSON(wr, http.StatusOK, resp)
}
