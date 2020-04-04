package db

func GetSlot(db *DB, slotID int) (Slot, error) {
	stmt := db.prepared["company/slot/get"]
	slot := Slot{}
	err := stmt.Get(&slot, slotID)
	
	return slot, err
}