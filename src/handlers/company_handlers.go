package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api/company"
	"github.com/gin-gonic/gin"
)

func init_company_routes(r *gin.Engine) {
	r.POST("/company", company.add_company)
	r.GET("/company", company.get_company)
	r.PATCH("/company", company.update_company)

	// to log in
	r.POST("/company/login", company.company_login)

	// with company token
	r.POST("/company/slots", company.add_slots)
	r.GET("/company/slots", company.get_slots)
	r.PATCH("/company/slots", company.update_slots)

	// get specific slot
	r.GET("/company/slots/id", company.get_slot) // /company/slots/id?id={id} -> c.Query("id")

	// scan qr code
	r.GET("/company/code", company.get_code) // /company/slots/code?code={code} -> c.Query("code")
}