package db

import (
	"database/sql"
	// "fmt"
	"time"
)

// reee ingen inbyggd insert-funktion
func InsertBooking(db *DB, book Booking) error {
	stmt := db.prepared["book/add"]
	_, err := stmt.Exec(book.SlotID, book.PhoneNumber, book.Code, book.FirstName, book.LastName, book.Visitee, book.Message, book.Status)

	return err
}

func GetBooking(db *DB, code string) (Booking, error) {
	stmt := db.prepared["book/get"]
	book := Booking{}
	err := stmt.Get(&book, code)

	return book, err
}

func RemoveBooking(db *DB, code string) error {
	stmt := db.prepared["book/remove"]
	_, err := stmt.Exec(code)

	return err
}

func GetBookingsBySlotID(db *DB, slotID int) ([]Booking, error) {
	stmt := db.prepared["booking/getBySlotID"]
	bookings := []Booking{}
	err := stmt.Select(&bookings, slotID)

	return bookings, err
}

func GetBookingsByCompanyID(db *DB, companyID int, startTime time.Time, endTime time.Time) (ret []BookingSlot, err error) {
	stmt := db.prepared["booking/getByCompanyID"]
	err = stmt.Select(&ret, companyID, startTime, endTime)
	return
}

func UpdateBookingStatus(db *DB, companyID int, bookingID int, status string) ([]Booking, error) {
	stmt := db.prepared["booking/update/status"]
	booking := []Booking{}
	err := stmt.Select(&booking, companyID, bookingID, status)

	if err == sql.ErrNoRows {
		return booking, nil
	}

	return booking, err
}

func UpdateBookingStatusByCode(db *DB, code string, status string) (ret Booking, err error) {
	stmt := db.prepared["book/code/update/status"]
	err = stmt.Get(&ret, code, status)
	return
}
