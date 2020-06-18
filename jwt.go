package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

//Holds the claims for the JWT.
//Tokens issued by this system will have the following claims
//Username: username of the user issued the token
//IssuedAt: time the token was issued at
//ExpiresAt: time the token will expire at
type jwtClaims struct {
	Username       string `json:"Username"`
	StandardClaims jwt.StandardClaims
}

var jwtSigningKey = []byte("Secret Keys In Code Are Bad")

func (token jwtClaims) Valid() error {
	if token.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("expired token")
	}
	return nil
}
