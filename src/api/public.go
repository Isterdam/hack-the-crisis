package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"

	"fmt"
	"encoding/json"
	"time"
	"math/rand"
	"strconv"
	"strings"
)

// book a time according to this JSON structure
type Booking struct {
	ID          int    `json:"id"`
	SlotID      int    `json:"slot_id"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

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
	var booking Booking
	err := json.NewDecoder(c.Request.Body).Decode(&booking)
	// could not parse enough arguments
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return
	}

	ticketCode := generateTicketCode(booking)
	// whitelist ticked code - to be checked at confirmation if it is contained
	Confirmed[ticketCode] = true

	// DOES THIS ACTUALLY GIVE CURRENT URL?
	url := c.Request.URL.RequestURI() + "/confirm/?code=" + ticketCode

	// generate confirmation url
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)
	
	store, _ := db.GetCompanyByID(dbbb, booking.ID)
	// TIME DOES NOT EXIST YET

	confirmation := "Hej " + booking.FirstName + "!\n\n" + "Vänligen bekräfta din bokning på " + store.Name.String + " klockan PLACEHOLDER"  + "\n\n" + url

	fmt.Println(confirmation)
	
	// UNCOMMENT THIS TO SEND TEXT
	// Send_text(c, booking.PhoneNum, confirmation)
}

func generateTicketCode(booking Booking) string {
	// current time + random num [0, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100)) + strings.ReplaceAll(booking.FirstName, " ", "_")
}

func Book_confirm(c *gin.Context) {
	ticketCode := c.Query("code")
	if Confirmed[ticketCode] {
		// add ticket to database here - THIS MAP WOULD OVERFLOW FROM UNCONFIRMED BOOKINGS, SHOULD BE EMPTIED FROM TIME TO TIME
		c.JSON(200, gin.H{
			"message": "Ticket confirmed yihoo!",
			"code": ticketCode,
		})
		delete(Confirmed, ticketCode) // delete entry from map
	} else {
		fmt.Println("Failed to verify ticket code")
	}
}

/*
func Unbook(c *gin.Context) {

}

func Get_ticket(c *gin.Context) {

}
*/
