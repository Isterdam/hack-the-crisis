package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
)

/*
func Add_company(c *gin.Context) {

}*/

func Get_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	comp, _ := db.GetCompanies(dbbb)

	c.JSON(200, comp)
}

/*
func Update_company(c *gin.Context) {

}
*/

func Add_slots(c *gin.Context) {
	if Is_authorized(c) {
		c.JSON(200, gin.H{
			"message": "Added slots successfully!",
		})
	}
}

func Get_slots(c *gin.Context) {
	if Is_authorized(c) {
		c.JSON(200, gin.H{
			"message": "Got slots successfully!",
		})
	}
}

func Update_slots(c *gin.Context) {
	if Is_authorized(c) {
		c.JSON(200, gin.H{
			"message": "Updated slots successfully!",
		})
	}
}

/*
func Get_slot(c *gin.Context) {

}

func Get_code(c *gin.Context) {

}
*/
