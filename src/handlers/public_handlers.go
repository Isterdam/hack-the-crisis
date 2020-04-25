package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitPublicRoutes(r *gin.Engine) {
	r.POST("/store/:store/slots", api.GetStoreSlots)
	r.GET("/slot/:slotID/load", api.GetSlotLoad) // amount booked and max number

	r.POST("/availability", api.GetCompanyAvailability)
	r.POST("/book", api.BookTime)                              // by phone number
	r.POST("/book/confirm/:code", api.ConfirmBookAndGetTicket) // by code
	r.POST("/unbook", api.Unbook)
}
