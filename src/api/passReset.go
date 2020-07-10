package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"
	"net/http"
	"os"
	"time"
)

// TODO: Invalidate previous tokens when a new password is set

type PassClaims struct {
	ID int
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
		int(comp.ID.Int64),
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

	content := "Hej!\n\nVi har mottagit din begäran om ett nytt lösenord på Booklie.se\n\nVänligen följ länken nedan för att skapa ett nytt lösenord\n\n" + url + "\n\nVänliga hälsningar,\nTeam Booklie"
	go SendMail(comp.Email.String, "Ditt lösenord på Booklie.se", content)

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
		return []byte(os.Getenv("JWTKEY")), nil
	})

	if payload, ok := token.Claims.(*PassClaims); ok && token.Valid {
		dbb := c.MustGet("db").(*db.DB)

		comp, err := db.GetCompanyByID(dbb, payload.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong!",
			})
			return
		}

		email = comp.Email.String

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

		content := "Hej!\n\nVi bekräftar härmed att ditt lösenord på Booklie.se har uppdaterats\n\nVänliga hälsningar,\nTeam Booklie"
		go SendMail(comp.Email.String, "Ditt lösenord har återställts", content)

		delete(ConfirmedHashes, email)

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
		return
	}
}
