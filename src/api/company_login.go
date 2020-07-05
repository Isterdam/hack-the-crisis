package api

import (
	"fmt"
	"os"

	"github.com/Isterdam/hack-the-crisis-backend/src/auth"
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"
	"net/http"
	"time"
)

// reads credentials from request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// CompanyLogin godoc
// @Summary Takes in a company as parameter, looks for password hash in database.
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Success 200 "Returns "success" and token as body. Also sets cookie with token."
// @Failure 401 "Unauthorized"
// @Router /company/login [post]
func CompanyLogin(c *gin.Context) {
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

	loginComp, err := db.GetCompanyByEmail(dbbb, comp.Email.String)

	if err != nil {
		fmt.Println(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginComp.Password.String), []byte(comp.Password.String))
	if err != nil {
		fmt.Println(err)
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	loginComp, err = db.GetCompanyByID(dbbb, int(loginComp.ID.Int64))

	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	// set expiration to now + 1 week
	// fill claims with username and standard
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not retrieve time zone!",
		})
		return
	}
	expirationTime := time.Now().In(loc).Add(168 * time.Hour)
	claims := &auth.Claims{
		ID: int(loginComp.ID.Int64),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("JWTKEY"))

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
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
		"token":   tokenString,
	})
}
