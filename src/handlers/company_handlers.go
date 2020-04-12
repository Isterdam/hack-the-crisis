package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/Isterdam/hack-the-crisis-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitCompanyRoutes(r *gin.Engine) {

	r.POST("/company", api.AddCompany)
	r.GET("/company", api.GetCompany)
	r.POST("/company/distance", api.GetCompanyDistance)

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
	}
	// get specific slot
	r.GET("/company/slots/id", api.GetSlot)

	// scan qr code
	r.POST("/company/code/:code/verify", api.VerifyCode)
	// check that booking is at the logged in company
}
