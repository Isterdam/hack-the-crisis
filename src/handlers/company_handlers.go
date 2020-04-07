package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitCompanyRoutes(r *gin.Engine) {
	r.PATCH("/company", api.UpdateCompany)

	r.POST("/company", api.AddCompany)
	r.GET("/company", api.GetCompany)
	r.POST("/company/distance", api.GetCompanyDistance)

	r.GET("/company/info", api.AuthGetCompany)

	// to log in
	r.POST("/company/login", api.CompanyLogin)

	// with company token
	r.POST("/company/slots", api.AddSlots)
	r.GET("/company/slots", api.GetSlots)
	r.PATCH("/company/slots", api.UpdateSlot)

	// get specific slot
	r.GET("/company/slots/id", api.GetSlot)

	// scan qr code
	r.POST("/company/code/:code/verify", api.VerifyCode)
	// check that booking is at the logged in company
}
