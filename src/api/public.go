package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"gopkg.in/guregu/null.v3"
	// "github.com/skip2/go-qrcode"
	"github.com/gin-gonic/gin"

	"fmt"
	"encoding/json"
	"time"
	"math/rand"
	"strconv"
	"strings"
)

/*
func Reserve_time(c *gin.Context) {

}

func Get_stores(c *gin.Context) {

}

func Get_store_slots(c *gin.Context) {

}

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

	// generate confirmation url
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)
	
	store, _ := db.GetCompanyByID(dbbb, int(booking.ID.Int64))
	// TIME DOES NOT EXIST YET

	confirmation := "Hej " + booking.FirstName.String + "!\n\n" + "Vänligen bekräfta din bokning på " + store.Name.String + " klockan PLACEHOLDER"  + "\n\n" + url

	fmt.Println(confirmation)
	
	// UNCOMMENT THIS TO SEND TEXT
	// Send_text(c, booking.PhoneNum, confirmation)
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

		url := strings.Split(c.Request.URL.Hostname() + c.Request.URL.Path, "?")[0] + "/get?code=" + ticketCode
		confirmation := "Du har nu bekräftat din bokning!\n\nBiljetten hittar du i länken nedan:\n" + url

		fmt.Println(confirmation)
		// Send_text(c, Confirmed[ticketCode].PhoneNumber.String, confirmation)

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
