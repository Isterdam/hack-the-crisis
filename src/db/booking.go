package db

// reee ingen inbyggd insert-funktion
func InsertBooking(db *DB, book Booking) error {
	stmt := db.prepared["book/add"]
	_, err := stmt.Exec(book.SlotID, book.PhoneNumber, book.Code, book.FirstName, book.LastName)
	
	return err
}

func GetBooking(db *DB, code string) (Booking, error) {
	stmt := db.prepared["book/get"]
	book := Booking{}
	err := stmt.Get(&book, code)
	
	return book, err
}

func RemoveBooking(db *DB, code string) (error) {
	stmt := db.prepared["book/remove"]
	_, err := stmt.Exec(code)

	return err
}