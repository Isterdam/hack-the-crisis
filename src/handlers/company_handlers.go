package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitCompanyRoutes(r *gin.Engine) {

	r.PATCH("/company", api.Update_company)

	r.POST("/company", api.Add_company)
	r.GET("/company", api.Get_company)
	r.POST("/company/distance", api.GetCompanyDistance)

	r.GET("/company/info", api.AuthGetCompany)

	// to log in
	r.POST("/company/login", api.CompanyLogin)

	// with company token
	r.POST("/company/slots", api.Add_slots)
	r.GET("/company/slots", api.Get_slots)
	r.PATCH("/company/slots", api.Update_slot)

	// get specific slot
	r.GET("/company/slots/id", api.Get_slot)

	// scan qr code
	r.POST("/company/code/:code/verify", api.VerifyCode)
	// check that booking is at the logged in company
}
