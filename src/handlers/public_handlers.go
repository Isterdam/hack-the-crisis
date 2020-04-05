package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/gin-gonic/gin"
)

func Init_public_routes(r *gin.Engine) {
	/*
		r.POST("/reserveBook", api.Reserve_time)

		r.GET("/stores", api.Get_stores) // by location and radius
	*/

	r.GET("/stores/:store/day/:day/slots", api.Get_store_slots)
	r.GET("/slot/:slotID/load", api.GetSlotLoad) // amount booked and max number

	/*
		r.GET("/search", api.Search_stores) // /search?word1={word1}&... -> c.Query("word1")
	*/

	r.POST("/book", api.Book_time)             // by phone number
	r.POST("/book/confirm", api.Book_confirm)  // by code
	r.POST("/unbook", api.Unbook)              // /unbook?code={code} -> c.Query("code")
	r.GET("/book/confirm/get", api.Get_ticket) // /book?code={code} -> c.Query("code")
}
