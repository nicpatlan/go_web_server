package main

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

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
