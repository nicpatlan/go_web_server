package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateUserToken(id, expiration int, jwtSecret string) (string, error) {
	defaultExpiration := 60 * 60 * 24
	if expiration == 0 {
		expiration = defaultExpiration
	} else if expiration > defaultExpiration {
		expiration = defaultExpiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "go_web_app",
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expiration))),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   strconv.Itoa(id),
	})
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenStr, jwtSecret string) (int, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return -1, err
	}
	id, err := token.Claims.GetSubject()
	if err == nil {
		id, err := strconv.Atoi(id)
		if err != nil {
			return -1, err
		}
		return id, nil
	}
	return -1, errors.New("could not validate token")
}
