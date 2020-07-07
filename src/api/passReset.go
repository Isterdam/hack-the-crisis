package api

import (
	"github.com/gin-gonic/gin"
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"

	"encoding/json"
	"fmt"
	"os"
	"time"
	"net/http"
)

// TODO: Invalidate previous tokens when a new password is set

type PassClaims struct {
	Email string 
	jwt.StandardClaims
}

func PasswordReset(c *gin.Context) {
	var comp db.Company

	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Body could not be parsed into Company struct!",
		})
		return
	}

	dbb := c.MustGet("db").(*db.DB)
	
	comp, err = db.GetCompanyByEmail(dbb, comp.Email.String)
	if err != nil {
		// company does not exist (probably - could be another error), but do not reveal it
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
		return
	}

	jwtKey := []byte(os.Getenv("JWTKEY"))

	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve time zone!",
		})
		return
	}
	expirationTime := time.Now().In(loc).Add(30 * time.Minute)

	passClaims := PassClaims{
		comp.Email.String, 
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, passClaims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong when trying to reset password!",
		})
		return
	}

	url := "https://www.booklie.se/password/reset/confirm/" + tokenString

	// whitelist hash
	ConfirmedHashes[comp.Email.String] = tokenString

	// send e-mail with confirmation link here
	fmt.Println(url)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	return
}

// both parameter for token and company body for new password
func PasswordResetToken(c *gin.Context) {
	tokenString := c.Param("token")
	var comp db.Company
	var email string

	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Body could not be parsed into company struct!",
		})
		return
	}

	pass := comp.Password.String

	token, _ := jwt.ParseWithClaims(tokenString, &PassClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("JWTKEY"), nil
	})
	
	if payload, ok := token.Claims.(*PassClaims); ok {
		email = payload.Email

		dbb := c.MustGet("db").(*db.DB)

		comp, err := db.GetCompanyByEmail(dbb, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong!",
			})
			return
		}

		if ConfirmedHashes[email] == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Not a valid token!",
			})
			return
		}
		
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Could not generate hash for password!",
			})
		}
	
		comp.Password.String = string(hash)

		_, err = db.UpdateCompany(dbb, comp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not update company in database!",
			})
			return
		}

		// Send email saying password was successfully changed here
		fmt.Println("Password was changed successfully!")

		delete(ConfirmedHashes, email)
	}
}