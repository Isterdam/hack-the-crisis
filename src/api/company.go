package api

import (
	"encoding/json"
	"fmt"
	"html"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/Isterdam/hack-the-crisis-backend/src/tz"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	null "gopkg.in/guregu/null.v3"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Company body could not be parsed correctly!",
			"error":   err.Error(),
		})
	}

	comp.Sanitize()

	if len(comp.Password.String) < 8 {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Password too short",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(comp.Password.String), bcrypt.MinCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not generate hash for password!",
			"error":   err.Error(),
		})
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
	for i, d := range ConfirmedCompanies {
		fmt.Printf("%#v %#v", i, d)
	}

	url := "www.shopalone.se" + c.Request.URL.Path + "/confirm/" + code
	msg := "Hello " + comp.Name.String + "!\n\n" + "Please confirm your email at ShopAlone in the link below:\n\n" + url

	// slow af so parallellize that shit
	go SendMail(comp.Email.String, "Confirm your email at ShopAlone", msg)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
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
	code = html.EscapeString(code)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)
	fmt.Println(code)
	for i, d := range ConfirmedCompanies {
		fmt.Printf("%#v %#v", i, d)
	}
	//fmt.Printf("%#v\n", ConfirmedCompanies[code])

	if ConfirmedCompanies[code].Email.String == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "The company does not exist yet!",
		})
	} else {
		// add company verified = true here?
		err := db.InsertCompany(dbbb, ConfirmedCompanies[code])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Company could not be inserted into database!",
				"error":   err.Error(),
			})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get companies from database!",
			"error":   err.Error(),
		})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Company body could not be parsed correctly!",
			"error":   err.Error(),
		})
	}

	comp.Sanitize()

	var newComp db.Company
	newComp, err = db.UpdateCompany(dbbb, comp)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not update the company!",
			"error":   err.Error(),
		})
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

	dbb := c.MustGet("db").(*db.DB)
	id := c.MustGet("id").(int)

	//var slots []db.Slot

	var req struct {
		Slots       []db.Slot `json:"slots"`
		Repetitions null.Int  `json:"repeat"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse json correctly.",
			"error":   err.Error(),
		})
		return
	}

	if req.Repetitions.Valid && req.Repetitions.Int64 < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Number of repetitions cannot be less than 0.",
		})
		return
	} else if !req.Repetitions.Valid {
		req.Repetitions = null.IntFrom(1)
	}

	for i := int64(0); i <= req.Repetitions.Int64; i++ {
		for _, slot := range req.Slots {
			slot.CompanyID = null.IntFrom(int64(id))
			slot.StartTime = null.TimeFrom(slot.StartTime.Time.AddDate(0, 0, int(i*7)))
			slot.EndTime = null.TimeFrom(slot.EndTime.Time.AddDate(0, 0, int(i*7)))

			err := db.AddSlot(dbb, slot)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Something went wrong!",
					"error":   err.Error(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get slots from database!",
			"error":   err.Error(),
		})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Slot body could not be parsed correctly!",
			"error":   err.Error(),
		})
	}

	var newSlot db.Slot
	newSlot, err = db.UpdateSlot(dbbb, slot)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not update slot in database!",
			"error":   err.Error(),
		})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Slot body could not be parsed correctly!",
			"error":   err.Error(),
		})
	}

	var newSlot db.Slot
	newSlot, err = db.GetSlot(dbbb, int(slot.ID.Int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get slot from database!",
			"error":   err.Error(),
		})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Distance body could not be parsed correctly!",
			"error":   err.Error(),
		})
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
		comps[i].DistToUser = null.FloatFrom(distance(dist.Latitude*(180/math.Pi), dist.Longitude*(180/math.Pi), float64(comps[i].Latitude.Float64), float64(comps[i].Longitude.Float64)))
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get companies from database!",
			"error":   err.Error(),
		})
	}

	c.JSON(200, comps)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get company from database!",
			"error":   err.Error(),
		})
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
	code = html.EscapeString(code)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Company body could not be parsed correctly!",
			"error":   err.Error(),
		})
		return
	}

	comp.Sanitize()

	loggedInCompanyID := int(comp.ID.Int64)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	booking, err := db.GetBooking(dbbb, code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get booking from database!",
			"error":   err.Error(),
		})
	}

	slot, err := db.GetSlot(dbbb, int(booking.SlotID.Int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not get slot from database!",
			"error":   err.Error(),
		})
	}

	bookingCompanyID := int(slot.CompanyID.Int64)

	loc, err := time.LoadLocation(tz.GetCountry(comp.Country.String).Zones[0].Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find the location for time zone!",
		})
	}
	t := time.Now().In(loc)

	validTime := t.After(slot.StartTime.Time) && t.Before(slot.EndTime.Time)
	if loggedInCompanyID == bookingCompanyID && validTime {
		c.JSON(200, gin.H{
			"message": "Ticket verified!",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ticket was not verified!",
		})
	}
}

func AddBookingAsCompany(c *gin.Context) {
	dbb := c.MustGet("db")
	dbbb := dbb.(*db.DB)

	var booking db.Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Booking body could not be parsed correctly!",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	booking.Code = null.StringFrom(ticketCode)

	booking.Sanitize()

	err = db.InsertBooking(dbbb, booking)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert booking into database!",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Booking was added successfully!",
		"data":    booking,
	})
}

func GetAllCompanyBookings(c *gin.Context) {
	dbbb := c.MustGet("db").(*db.DB)
	id := c.MustGet("id").(int)

	startTime := c.Query("start")
	endTime := c.Query("end")

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "The query for start time is not a valid timestamp!" +
				" The correct layout is e.g 2020-06-28T15:11:50.341Z (in UTC)",
			"layout": "Your layout was " + startTime,
		})
		return
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "The query for end time is not a valid timestamp!" +
				" The correct layout is e.g 2020-06-28T15:11:50.341Z (in UTC)",
			"layout": "Your layout was " + endTime,
		})
		return
	}

	bookings, err := db.GetBookingsByCompanyID(dbbb, id, start, end)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get bookings from database!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    bookings,
	})

}

func UpdateCompanyBookingStatus(c *gin.Context) {
	var req struct {
		Status string
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	req.Status = html.EscapeString(req.Status)

	dbb := c.MustGet("db")
	dbbb := dbb.(*db.DB)

	id := c.MustGet("id")

	bookingID, err := strconv.Atoi(c.Param("bookingID"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID",
		})
		return
	}

	updatedBooking := []db.Booking{}
	updatedBooking, err = db.UpdateBookingStatus(dbbb, id.(int), bookingID, req.Status)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Database Error",
		})
		return
	}

	if len(updatedBooking) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "This booking does not exist",
		})
	} else {
		if req.Status == "cancelled" {
			slot, err := db.GetSlot(dbbb, int(updatedBooking[0].SlotID.Int64))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Could not get the time slot from database!",
				})
			}

			stockholmTime, err := time.LoadLocation("Europe/Stockholm") // sorry this is hard-coded
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Could not get time zone!",
				})
			}

			start := slot.StartTime.Time.In(stockholmTime)
			stop := slot.EndTime.Time.In(stockholmTime)

			text := "Hej " + updatedBooking[0].FirstName.String + "!\n\nDin bokning den " + start.Format("2/1") + " klockan " + start.Format("15:04") + "-" + stop.Format("15:04") + " har tyvärr ställts in.\n\nVi beklagar detta!"

			go Send_text(c, updatedBooking[0].PhoneNumber.String, text)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    updatedBooking[0], //The array will always contain only one booking
		})
	}

}
