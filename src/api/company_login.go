package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"time"
)

// jwt key used to create signature
var jwtKey = []byte(JWTkey)

// reads credentials from request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// will be encoded to JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Company_login(c *gin.Context) {
	var comp db.Company
	// decode HTTP Request Body JSON
	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return
	}

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	id, err := db.VerifyLoginCompany(dbbb, comp)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	loginComp, err := db.GetCompanyByID(dbbb, int(id))

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	// set expiration to now + 30 mins
	// fill claims with username and standard
	loc, _ := time.LoadLocation("Europe/Stockholm")
	expirationTime := time.Now().In(loc).Add(12 * time.Hour)
	claims := &Claims{
		Email: loginComp.Email.String,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}

	// set cookie with token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	c.JSON(200, gin.H{
		"message": "Success",
		"token": tokenString,
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
