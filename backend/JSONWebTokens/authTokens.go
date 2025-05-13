package authToken

import "github.com/golang-jwt/jwt/v4"

type JWTToken struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.RegisteredClaims
}
