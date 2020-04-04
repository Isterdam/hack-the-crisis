package handlers

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/api"
	"github.com/gin-gonic/gin"
)

func Init_company_routes(r *gin.Engine) {

	r.PATCH("/company", api.Update_company)

	r.POST("/company", api.Add_company)
	r.GET("/company", api.Get_company)

	// to log in
	r.POST("/company/login", api.Company_login)

	// with company token
	r.POST("/company/slots", api.Add_slots)
	r.GET("/company/slots", api.Get_slots)
	r.PATCH("/company/slots", api.Update_slots)

	/*
		// get specific slot
		r.GET("/company/slots/id", api.Get_slot) // /company/slots/id?id={id} -> c.Query("id")

		// scan qr code
		r.GET("/company/code", api.Get_code) // /company/slots/code?code={code} -> c.Query("code")
	*/
}
