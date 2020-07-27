package api

import (
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/Isterdam/hack-the-crisis-backend/src/utils/random"
	null "gopkg.in/guregu/null.v3"

	"github.com/gin-gonic/gin"
)

func CreateBooking(c *gin.Context) {
	dbb := c.MustGet("db").(*db.DB)
	var booking db.Booking

	err := c.ShouldBindJSON(&booking)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON.",
		})
	}

	booking.Code = null.StringFrom(utils.RandStringBytes(6))

	newBooking, err := db.InsertBooking(dbb, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert booking.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    newBooking,
	})
}
