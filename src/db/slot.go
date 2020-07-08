package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func GetSlot(db *DB, slotID int) (Slot, error) {
	stmt := db.prepared["company/slot/get"]
	slot := Slot{}
	err := stmt.Get(&slot, slotID)

	return slot, err
}

func GetSlotsByCompany(db *DB, companyID int) ([]Slot, error) {
	stmt := db.prepared["company/slot/getAll"]
	slots := []Slot{}
	err := stmt.Select(&slots, companyID)

	return slots, err
}

func GetSlotsByCompanyAndBetween(db *DB, companyID int, start time.Time, end time.Time) ([]Slot, error) {
	stmt := db.prepared["company/slot/get/betweenTime"]
	slots := []Slot{}
	err := stmt.Select(&slots, start, end, companyID)
	return slots, err
}

func AddSlot(db *DB, slot Slot) error {
	stmt := db.prepared["company/slot/add"]
	_, err := stmt.Exec(slot.CompanyID, slot.StartTime, slot.EndTime, slot.MaxAmount)

	return err
}

func UpdateSlot(db *DB, slot Slot) (Slot, error) {
	var newSlot Slot

	stmt := db.prepared["company/slot/update"]
	err := stmt.QueryRowx(slot.ID, slot.StartTime, slot.EndTime, slot.MaxAmount, slot.Booked).StructScan(&newSlot)

	return newSlot, err
}

func DeleteSlots(db *DB, slotIDs []int) ([]Slot, error) {
	deletedSlots := []Slot{}
	query := "DELETE FROM slots WHERE id IN (?) RETURNING *"
	query, args, err := sqlx.In(query, slotIDs)

	query = db.DB.Rebind(query)
	stmt, err := db.DB.Preparex(query)

	rows, err := stmt.Queryx(args...)

	for rows.Next() {
		var slot Slot
		rows.StructScan(&slot)
		deletedSlots = append(deletedSlots, slot)
	}

	return deletedSlots, err
}

func GetSlotsByID(db *DB, slotIDs []int, cID int) ([]Slot, error) {
	var slots []Slot
	query := "SELECT * from slots WHERE company_id = ? AND id IN (?)"
	newq, args, err := sqlx.In(query, cID, slotIDs)

	query = db.DB.Rebind(newq)
	stmt, err := db.DB.Preparex(query)
	err = stmt.Select(&slots, args...)

	return slots, err
}

func UpdateSlotBooked(db *DB, slotID int, booked int) (slot Slot, err error) {
	stmt := db.prepared["company/slot/update/booked"]
	err = stmt.Get(&slot, slotID, booked)
	return
}
