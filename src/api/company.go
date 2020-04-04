package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Add_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}

	dbbb := dbb.(*db.DB)
	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	if err != nil {
		fmt.Printf("hello2 %s", err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(comp.Password.String), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	comp.Password.String = string(hash)

	err = db.InsertCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
	})
}

func Get_company(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	comp, _ := db.GetCompanies(dbbb)

	c.JSON(200, comp)
}

func Update_company(c *gin.Context) {
	if !Is_authorized(c) {
		return
	}
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var comp db.Company
	err := json.NewDecoder(c.Request.Body).Decode(&comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	var newComp db.Company
	newComp, err = db.UpdateCompany(dbbb, comp)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	c.JSON(200, newComp)
}

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
