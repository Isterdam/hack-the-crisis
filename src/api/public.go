package api

import (
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v3"

	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetStoreSlots godoc
// @Summary Gets all slots for a certain company on a certain day.
// @Produce json
// @Param day path string true "Day"
// @Param store path string true "Store"
// @Success 200 {array} db.Slot
// @Router /stores/{store}/day/{day}/slots [get]
func GetStoreSlots(c *gin.Context) {
	dayStr := c.Param("day")
	day, _ := strconv.Atoi(dayStr)

	storeIDStr := c.Param("store")
	storeID, _ := strconv.Atoi(storeIDStr)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	slots, _ := db.GetSlotsByCompany(dbbb, storeID)

	var slotsByDay []db.Slot
	for _, slot := range slots {
		if int(slot.Day.Int64) == day {
			slotsByDay = append(slotsByDay, slot)
		}
	}

	c.JSON(http.StatusOK, slotsByDay)
}

// BookTime godoc
// @Summary "Books" a certain time by creating a confirmation link that is sent to the user by text. Does NOT add booking to database.
// @Consume json
// @Param booking body db.Booking true "Booking"
// @Router /book [post]
func BookTime(c *gin.Context) {
	var booking db.Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	// could not parse enough arguments
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Page not found",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	booking.Code = null.StringFrom(ticketCode)
	// whitelist ticked code - to be checked at confirmation if it is contained
	Confirmed[ticketCode] = booking

	url := "www.shopalone.se" + c.Request.URL.Path + "/confirm/" + ticketCode

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	timeSlot, _ := db.GetSlot(dbbb, int(booking.SlotID.Int64))
	store, _ := db.GetCompanyByID(dbbb, int(timeSlot.CompanyID.Int64))

	timeStart := timeSlot.StartTime.Time
	timeStop := timeSlot.EndTime.Time

	var timeStr string
	timeStr = timeStart.Format("15:04") + "-" + timeStop.Format("15:04")

	confirmation := "Hej " + booking.FirstName.String + "!\n\n" + "Vänligen bekräfta din bokning på " + store.Name.String + " klockan " + timeStr + " och hämta sedan din biljett genom länken nedan:\n\n" + url

	fmt.Println(confirmation)
	// Send_text(c, booking.PhoneNumber.String, confirmation)

	c.JSON(200, gin.H{
		"message": "Booking was successful",
	})
}

func generateTicketCode(booking db.Booking) string {
	// last 2 digits of current time + random num [10, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10)[8:] + strconv.Itoa(10+rand.Intn(90)) + strings.ReplaceAll(booking.FirstName.String, " ", "_")
}

// ConfirmBookAndGetTicket godoc
// @Summary Confirms a booking and adds it to the database if first time. Gets a ticket if it has already been added to database.
// @Produce json
// @Param code path string true "Code"
// @Router /book/confirm/{code} [post]
func ConfirmBookAndGetTicket(c *gin.Context) {
	code := c.Param("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	booking, _ := db.GetBooking(dbbb, code)

	if Confirmed[code].PhoneNumber.String == "" && booking.PhoneNumber.String == "" {
		// booking does not exist
		fmt.Println("Booking does not exist!")
		return
	} else if Confirmed[code].PhoneNumber.String == "" && booking.PhoneNumber.String != "" {
		// booking exists and has been added to database
		c.JSON(200, booking)
	} else {
		// booking exists but has not yet been added to database
		db.InsertBooking(dbbb, Confirmed[code])
		booking, _ = db.GetBooking(dbbb, code)
		c.JSON(200, booking)
		delete(Confirmed, code) // delete entry from map
	}
}

// Unbook godoc
// @Summary Unbooks a ticket by removing it from the database by code.
// @Param code path string true "Code"
// @Router /unbook [post]
func Unbook(c *gin.Context) {
	code := c.Param("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	err := db.RemoveBooking(dbbb, code)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not remove booking!")
	}
}

// GetSlotLoad godoc
// @Summary Gets the load of a slot by returning maxAmount of customers and amount of booked customers as JSON.
// @Produce json
// @Param slotID path string true "slotID"
// @Success 200 "JSON with "maxAmount", "bookingsAmount""
// @Router /slot/{slotID}/load [get]
func GetSlotLoad(c *gin.Context) {
	slotIDStr := c.Param("slotID")
	slotID, _ := strconv.Atoi(slotIDStr)

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	slot, err := db.GetSlot(dbbb, slotID)
	if err != nil {
		fmt.Println(err)
		return
	}

	maxAmount := strconv.Itoa(int(slot.MaxAmount.Int64))

	bookings, err := db.GetBookingsBySlotID(dbbb, slotID)
	if err != nil {
		fmt.Println(err)
		return
	}

	bookingsAmount := strconv.Itoa(len(bookings))

	c.JSON(200, gin.H{
		"maxAmount":      maxAmount,
		"bookingsAmount": bookingsAmount,
	})
}
