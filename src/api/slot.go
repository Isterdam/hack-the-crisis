package api

import (
	"encoding/json"
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
)

// DeleteSlots godoc
// @Summary Takes in JSON array of slot IDs, deletes them in database and returns the deleted slots
// @Consume json
// @Produce json
// @Param slotIDs body []integer true "Slot IDs"
// @Success 200 {array} db.Slot
// @Router /company/slots [delete]
func DeleteSlots(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	ID, exist := c.Get("id")

	if !exist {
		return
	}

	var slotIDs []int
	err := json.NewDecoder(c.Request.Body).Decode(&slotIDs)

	if err != nil {
		return
	}

	slots, err := db.GetSlotsByID(dbbb, slotIDs, ID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get slot from database by ID!",
		})
		return
	}

	if len(slots) != len(slotIDs) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Fail",
		})
		return
	}

	dSlots, err := db.DeleteSlots(dbbb, slotIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete slots from database by IDs!",
		})
		return
	}

	c.JSON(http.StatusOK, dSlots)
}
