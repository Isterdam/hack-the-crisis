package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AddCompany godoc
// @Summary Sends an email to company asking them to confirm
// @Consume json
// @Produce json
// @Param company body db.Company true "Company"
// @Success 200
// @Router /company [post]
func AddCompany(c *gin.Context) {
	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	fmt.Printf("%#v", comp)
	if err != nil {
		fmt.Printf("hello2 %s", err)
		return
	}

	if len(comp.Password.String) < 8 {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Password too short",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(comp.Password.String), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	comp.Password.String = string(hash)

	comp.Latitude.Float64 = comp.Latitude.Float64 / 180 * math.Pi
	comp.Longitude.Float64 = comp.Longitude.Float64 / 180 * math.Pi

	code := generateVerifyingCode(comp)
	r := strings.NewReplacer("ä", "a",
		"å", "a",
		"ö", "o")
	code = r.Replace(code)
	ConfirmedCompanies[code] = comp

	url := "www.shopalone.se" + c.Request.URL.Path + "/confirm/" + code
	msg := "Hello " + comp.Name.String + "!\n\n" + "Please confirm your email at ShopAlone in the link below:\n\n" + url

	// slow af so parallellize that shit
	go SendMail(comp.Email.String, "Confirm your email at ShopAlone", msg)
}

func generateVerifyingCode(company db.Company) string {
	// last 2 digits of current time + random num [10, 100) + company name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10)[8:] + strconv.Itoa(10+rand.Intn(90)) + strings.ReplaceAll(company.Name.String, " ", "_")
}

// ConfirmCompany godoc
// @Summary Confirms a company and adds it to the database
// @Produce json
// @Param code path string true "Code"
// @Success 200
// @Router /company/confirm/{code} [post]
func ConfirmCompany(c *gin.Context) {
	code := c.Param("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	if ConfirmedCompanies[code].Email.String == "" {
		fmt.Println("Company does not exist!")
		return
	} else {
		// add company verified = true here?
		err := db.InsertCompany(dbbb, ConfirmedCompanies[code])
		if err != nil {
			fmt.Println(err)
			return
		}
		delete(ConfirmedCompanies, code)
		c.JSON(200, gin.H{
			"message": "Company successfully added!",
		})
	}
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

	comp, err := db.GetCompaniesVerifiedPublic(dbbb)

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

	id, exist := c.Get("id")
	if !exist {
		return
	}

	var slots []db.Slot
	slots, err := db.GetSlotsByCompany(dbbb, id.(int))

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

	dist.Latitude = dist.Latitude / 180 * math.Pi
	dist.Longitude = dist.Longitude / 180 * math.Pi

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

	for i := range comps {
		comps[i].Latitude.Float64 = comps[i].Latitude.Float64 / math.Pi * 180
		comps[i].Longitude.Float64 = comps[i].Longitude.Float64 / math.Pi * 180
	}

	if err != nil {
		return
	}

	c.JSON(200, comps)
}

// SearchForCompanies godoc
// @Summary Iteratively searches for companies, starting with 5 kilometers and increasing, until at least 10 companies have been found
// @Produce json
// @Param lon path number true "Longitude"
// @Param lat path number true "Latitude"
// @Success 200 {array} db.CompanyPublic
// @Router /company/search/{lon}/{lat} [get]
func SearchForCompanies(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var dist db.Distance

	lon, _ := strconv.ParseFloat(c.Param("lon"), 64)
	lat, _ := strconv.ParseFloat(c.Param("lat"), 64)

	dist.Latitude = lat / 180 * math.Pi
	dist.Longitude = lon / 180 * math.Pi

	dist.Distance = 5
	var comps map[db.CompanyPublic]bool
	comps = make(map[db.CompanyPublic]bool)

	for i := 0; i < 6; i++ {
		if len(comps) >= 10 {
			break
		}
		dist.R = float64(dist.Distance) / 6371

		r := dist.R
		dist.LatMax = dist.Latitude + r
		dist.LatMin = dist.Latitude - r

		dlon := math.Asin(math.Sin(r) / math.Cos(dist.Latitude))

		dist.LonMin = dist.Longitude - dlon
		dist.LonMax = dist.Longitude + dlon

		compsAppend, err := db.GetCompaniesWithinDistance(dbbb, dist)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, comp := range compsAppend {
			comps[comp] = true
		}

		dist.Distance = int(float64(dist.Distance) * increaseDistance(len(comps)))
	}

	var compsSlice []db.CompanyPublic
	for comp, _ := range comps {
		compsSlice = append(compsSlice, comp)
	}

	c.JSON(200, compsSlice)
}

// super proprietary function to increase distance while searching for stores
// parameter x is #stores currently found
func increaseDistance(x int) float64 {
	return 1 + (2 / (math.Exp(float64(x))))
}

// AuthGetCompany godoc
// @Summary Gets a full company by id, no password required. Requires authorization. Gets company from context.
// @Consume json
// @Produce json
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

	validTime := time.Now().After(slot.StartTime.Time) && time.Now().Before(slot.EndTime.Time)
	if loggedInCompanyID == bookingCompanyID && validTime {
		c.JSON(200, gin.H{
			"message": "Ticket verified!",
		})
	} else {
		c.JSON(401, gin.H{
			"message": "Ticket was not verified!",
		})
	}
}
