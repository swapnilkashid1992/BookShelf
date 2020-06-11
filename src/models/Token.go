package models

import "github.com/dgrijalva/jwt-go"

type Token struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
