package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Add_company godoc
// @Summary Adds a company to the database
// @Consume json
// @Produce json
// @Param id path integer true "User ID"
// @Success 200
// @Router /company [post]
func Add_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}

	dbbb := dbb.(*db.DB)
	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	fmt.Printf("%#v", comp)
	if err != nil {
		fmt.Printf("hello2 %s", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(comp.Password.String), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	comp.Password.String = string(hash)

	err = db.InsertCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
	})
}

func Get_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	comp, err := db.GetCompaniesPublic(dbbb)

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, comp)
}

func Update_company(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	var newComp db.Company
	newComp, err = db.UpdateCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	c.JSON(200, newComp)
}

func Add_slots(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slots []db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slots)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, slot := range slots {
		db.AddSlot(dbbb, slot)
	}
}

func Get_slots(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var slots []db.Slot
	slots, err = db.GetSlotsByCompany(dbbb, int(comp.ID.Int64))
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, slots)
}

func Update_slot(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slot db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newSlot db.Slot
	newSlot, err = db.UpdateSlot(dbbb, slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, newSlot)
}

func Get_slot(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var slot db.Slot
	err := json.NewDecoder(c.Request.Body).Decode(&slot)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newSlot db.Slot
	newSlot, err = db.GetSlot(dbbb, int(slot.ID.Int64))
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, newSlot)
}

func GetCompanyDistance(c *gin.Context) {
	var dist db.Distance
	err := json.NewDecoder(c.Request.Body).Decode(&dist)
	if err != nil {
		fmt.Println(err)
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return
	}

	dist.R = float64(dist.Distance) / 6371

	r := dist.R
	dist.LatMax = dist.Latitude + r
	dist.LatMin = dist.Latitude - r

	dlon := math.Asin(math.Sin(r) / math.Cos(dist.Latitude))

	dist.LonMin = dist.Longitude - dlon
	dist.LonMax = dist.Longitude + dlon

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	comps, err := db.GetCompaniesWithinDistance(dbbb, dist)

	if err != nil {
		return
	}

	c.JSON(200, comps)
}

func AuthGetCompany(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)
	token := c.Request.Header.Get("Authorization")

	// jwt string from token
	claims := &Claims{}

	// parse jwt and store in claims
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	ccc := tkn.Claims.(*Claims)

	comp, err := db.GetCompanyByIDNoPass(dbbb, ccc.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, comp)
	return
}

func VerifyCode(c *gin.Context) {
	if !IsAuthorized(c) {
		return
	}
	code := c.Query("code")

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		fmt.Println(err)
		return
	}

	loggedInCompanyID := int(comp.ID.Int64)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	booking, err := db.GetBooking(dbbb, code)
	if err != nil {
		fmt.Println(err)
		return
	}

	slot, err := db.GetSlot(dbbb, int(booking.SlotID.Int64))
	if err != nil {
		fmt.Println(err)
		return
	}

	bookingCompanyID := int(slot.CompanyID.Int64)

	if loggedInCompanyID == bookingCompanyID {
		c.JSON(200, gin.H{
			"message": "Ticket verified!",
		})
	} else {
		c.JSON(401, gin.H{
			"message": "Ticket was not verified!",
		})
	}
}
