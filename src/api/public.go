package api

import (
	"github.com/gin-gonic/gin"

	"encoding/json"
)

// book a time according to this JSON structure
type Booking struct {
	Name string `json:"name"`
	Time string `json:"time"`
	Store string `json:"store"`
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

	// generate confirmation url here somewhere
	confirmation := "Hej "+booking.Name+"!\n\n"+"Vänligen bekräfta din bokning på "+booking.Store+" klockan "+booking.Time
	Send_text(c, booking.PhoneNum, confirmation)
}

/*
func Book_confirm(c *gin.Context) {

}

func Unbook(c *gin.Context) {

}

func Get_ticket(c *gin.Context) {

}
*/