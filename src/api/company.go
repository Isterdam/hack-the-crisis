package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AddCompany godoc
// @Summary Adds a company to the database
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Success 200
// @Router /company [post]
func AddCompany(c *gin.Context) {
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

// GetCompany godoc
// @Summary Gets all CompanyPublic from database
// @Produce json
// @Success 200 {array} db.CompanyPublic
// @Router /company [get]
func GetCompany(c *gin.Context) {
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

// UpdateCompany godoc
// @Summary Updates a company in the database, then returns the updated company. Requires authorization.
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Success 200 {object} db.Company
// @Router /company [patch]
func UpdateCompany(c *gin.Context) {
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

// AddSlots godoc
// @Summary Adds slots to database. Requires authorization.
// @Consume json
// @Produce json
// @Param slots body []db.Slot true "Slots"
// @Success 200
// @Router /company/slots [post]
func AddSlots(c *gin.Context) {

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

// GetSlots godoc
// @Summary Gets all slots for a certain company. Requires authorization.
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Success 200 {array} db.Slot
// @Router /company/slots [get]
func GetSlots(c *gin.Context) {

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

// UpdateSlot godoc
// @Summary Updates a certain slot, then returns updated slot. Requires authorization.
// @Consume json
// @Produce json
// @Param slot body db.Slot true "Slot"
// @Success 200 {object} db.Slot
// @Router /company/slots [patch]
func UpdateSlot(c *gin.Context) {

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

// GetSlot godoc
// @Summary Gets a full slot. Requires a slot as parameter, but an id in body will suffice.
// @Consume json
// @Produce json
// @Param slot body db.Slot true "Slot"
// @Success 200 {object} db.Slot
// @Router /company/slots/id [get]
func GetSlot(c *gin.Context) {
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

// GetCompanyDistance godoc
// @Summary Gets public companies within a certain distance.
// @Consume json
// @Produce json
// @Param distance body db.Distance true "Distance"
// @Success 200 {array} db.CompanyPublic
// @Router /company/distance [post]
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

// AuthGetCompany godoc
// @Summary Gets a full company by id, no password required. Requires authorization.
// @Consume json
// @Produce json
// @Param authorization header string true "Token"
// @Success 200 {object} db.Company
// @Router /company/info [get]
func AuthGetCompany(c *gin.Context) {

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	ID, exist := c.Get("id")

	if !exist {
		return
	}

	comp, err := db.GetCompanyByIDNoPass(dbbb, ID.(int))
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, comp)
	return
}

// VerifyCode godoc
// @Summary Verifies a ticket code for a company. Requires authorization.
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Param code path string true "Code"
// @Success 200 "Ticket was verified."
// @Failure 401 "Ticket could not be verified."
// @Router /company/code/{code}/verify [post]
func VerifyCode(c *gin.Context) {
	code := c.Param("code")
	fmt.Println(code)

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
