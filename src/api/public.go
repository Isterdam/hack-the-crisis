package api

import (
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
	Name     string `json:"name"`
	Time     string `json:"time"`
	Store    string `json:"store"`
	PhoneNum string `json:"phone_num"`
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
	// DOES THIS ACTUALLY GIVE CURRENT URL?
	url := c.Request.URL.RequestURI() + "/confirm/?code=" + ticketCode
	// generate confirmation url here somewhere
	confirmation := "Hej " + booking.Name + "!\n\n" + "Vänligen bekräfta din bokning på " + booking.Store + " klockan " + booking.Time + "\n\n" + url

	fmt.Println(confirmation)
	
	// UNCOMMENT THIS TO SEND TEXT
	// Send_text(c, booking.PhoneNum, confirmation)
}

func generateTicketCode(booking Booking) string {
	// current time + random num [0, 100) + booking name (where space is replaced by underscore)
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100)) + strings.ReplaceAll(booking.Name, " ", "_")
}

func Book_confirm(c *gin.Context) {
	ticketCode := c.Query("code")
	// add ticket to database here
	c.JSON(200, gin.H{
		"message": "Ticket confirmed yihoo!",
		"code": ticketCode,
	})
}

/*
func Unbook(c *gin.Context) {

}

func Get_ticket(c *gin.Context) {

}
*/
