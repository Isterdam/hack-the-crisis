package auth

import (
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// will be encoded to JWT
type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func IsValidToken(token string, c *Claims) bool {
	// obtain token from session cookies
	jwtKey := []byte(os.Getenv("JWTKEY"))

	if token == "" {
		return false
	}

	// jwt string from token
	claims := &Claims{}

	// parse jwt and store in claims
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {

		return false
	}
	if !tkn.Valid {

		return false
	}

	c = tkn.Claims.(*Claims)

	// token is valid
	return true
}
