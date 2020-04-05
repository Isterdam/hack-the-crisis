package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"gopkg.in/guregu/null.v3"
	// "github.com/skip2/go-qrcode"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*
func Reserve_time(c *gin.Context) {

}

func Get_stores(c *gin.Context) {

}
*/

func Get_store_slots(c *gin.Context) {
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

	c.JSON(200, slotsByDay)
}

/*
func Search_stores(c *gin.Context) {

}
*/

func Book_time(c *gin.Context) {
	var booking db.Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	// could not parse enough arguments
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	booking.Code = null.StringFrom(ticketCode)
	// whitelist ticked code - to be checked at confirmation if it is contained
	Confirmed[ticketCode] = booking

	// DOES THIS ACTUALLY GIVE CURRENT URL?
	url := c.Request.URL.Hostname() + c.Request.URL.Path + "/confirm/?code=" + ticketCode

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	store, _ := db.GetCompanyByID(dbbb, int(booking.ID.Int64))
	timeSlot, _ := db.GetSlot(dbbb, int(booking.SlotID.Int64))

	timeStart := timeSlot.StartTime.Time
	timeStop := timeSlot.EndTime.Time

	var timeStr string
	timeStr = strconv.Itoa(timeStart.Hour()) + ":" + strconv.Itoa(timeStart.Minute()) + "-" + strconv.Itoa(timeStop.Hour()) + ":" + strconv.Itoa(timeStop.Minute())

	confirmation := "Hej " + booking.FirstName.String + "!\n\n" + "Vänligen bekräfta din bokning på " + store.Name.String + " klockan " + timeStr + "\n\n" + url

	fmt.Println(confirmation)

	Send_text(c, booking.PhoneNumber.String, confirmation)
}

func generateTicketCode(booking db.Booking) string {
	// current time + random num [0, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100)) + strings.ReplaceAll(booking.FirstName.String, " ", "_")
}

func Book_confirm(c *gin.Context) {
	ticketCode := c.Query("code")
	if Confirmed[ticketCode].PhoneNumber.String != "" {
		dbb, exist := c.Get("db")
		if !exist {
			return
		}
		dbbb := dbb.(*db.DB)

		db.InsertBooking(dbbb, Confirmed[ticketCode])

		url := strings.Split(c.Request.URL.Hostname()+c.Request.URL.Path, "?")[0] + "/get?code=" + ticketCode
		confirmation := "Du har nu bekräftat din bokning!\n\nBiljetten hittar du i länken nedan:\n" + url

		fmt.Println(confirmation)
		Send_text(c, Confirmed[ticketCode].PhoneNumber.String, confirmation)

		// qrPng, _ := qrcode.Encode(url, qrcode.Medium, 256)
		// save qr code somewhere

		delete(Confirmed, ticketCode) // delete entry from map
	} else {
		fmt.Println("Failed to verify ticket code")
	}
}

func Unbook(c *gin.Context) {
	code := c.Query("code")

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

func Get_ticket(c *gin.Context) {
	code := c.Query("code")

	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	book, err := db.GetBooking(dbbb, code)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not find booking")
	}

	c.JSON(200, book)
}

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
