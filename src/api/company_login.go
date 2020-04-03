package api 

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

	"encoding/json"
	"net/http"
	"time"
)

// jwt key used to create signature
var jwtKey = []byte(JWTkey)

var users = map[string]string {
	"admin": "admin",
}

// reads credentials from request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// will be encoded to JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Company_login(c *gin.Context) {
	var creds Credentials 
	// decode HTTP Request Body JSON
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return
	}

	// GET PASSWORD FROM DATABASE OR SOMETHING HERE
	expectedPassword, ok := users[creds.Username] // placeholder

	if !ok || expectedPassword != creds.Password {
		c.JSON(404, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	// set expiration to now + 30 mins
	// fill claims with username and standard
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Internal server error",
		})
		return
	}

	// set cookie with token
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func Is_authorized(c *gin.Context) bool {
	// obtain token from session cookies
	t, err := c.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(404, gin.H{
				"message": "Unauthorized",
			})
			return false
		}
		c.JSON(404, gin.H{
			"message": "Bad request",
		})
		return false
	}

	// jwt string from token
	tknStr := t.Value
	claims := &Claims{}

	// parse jwt and store in claims
	tkn, err := jwt.ParseWithClaims(tknStr, claims, 
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(404, gin.H{
				"message": "Unauthorized",
			})
			return false
		}
		c.JSON(404, gin.H{
			"message": "Bad request",
		})
		return false
	}
	if !tkn.Valid {
		c.JSON(404, gin.H{
			"message": "Unauthorized",
		})
		return false
	}

	// token is valid
	return true
}