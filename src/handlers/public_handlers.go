package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api/"
	"github.com/gin-gonic/gin"
)

func init_public_routes(r *gin.Engine) {
	// r.POST("/reserveBook", public.reserve_time)

	r.GET("/stores", public.get_stores) // by location and radius

	/*
	r.GET("/stores/slots", public.get_store_slots) // by day

	r.GET("/search", public.search_stores) // /search?word1={word1}&... -> c.Query("word1")

	r.POST("/book", public.book_time) // by phone number
	r.POST("/book/confirm", public.book_confirm) // by code
	r.POST("/unbook", public.unbook)

	r.GET("/book", public.get_ticket) // /book?code={code} -> c.Query("code")
	*/
}