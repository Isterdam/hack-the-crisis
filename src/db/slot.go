package db

import (
	"time"
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
	err := stmt.QueryRowx(slot.ID, slot.StartTime, slot.EndTime, slot.MaxAmount).StructScan(&newSlot)

	return newSlot, err
}
