package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/Isterdam/hack-the-crisis-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitCompanyRoutes(r *gin.Engine) {

	r.POST("/company", api.AddCompany)
	r.POST("/company/confirm/:code", api.ConfirmCompany)
	r.GET("/company", api.GetCompany)
	r.POST("/company/distance", api.GetCompanyDistance)
	r.GET("/company/search/:lon/:lat", api.SearchForCompanies)

	// to log in
	r.POST("/company/login", api.CompanyLogin)

	// with company token

	auth := r.Group("/", middleware.AuthRequired)
	{
		auth.PATCH("/company", api.UpdateCompany)
		auth.GET("/company/info", api.AuthGetCompany)
		auth.POST("/company/slots", api.AddSlots)
		auth.GET("/company/slots", api.GetSlots)
		auth.PATCH("/company/slots", api.UpdateSlot)
		auth.DELETE("/company/slots", api.DeleteSlots)

		auth.GET("/company/booking", api.GetAllCompanyBookings)
		auth.PATCH("/company/booking/:bookingID/status", api.UpdateCompanyBookingStatus)
	}
	// get specific slot
	r.GET("/company/slots/id", api.GetSlot)

	// scan qr code
	r.POST("/company/code/:code/verify", api.VerifyCode)
	// check that booking is at the logged in company
}
