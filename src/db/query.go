package db

type KVP struct {
	K string
	V string
}

var queries = []KVP{
	KVP{K: "company/get", V: "SELECT * FROM company"},
	KVP{K: "company/getByID", V: "SELECT * FROM company WHERE id=$1"},
	KVP{K: "book/add", V: "INSERT INTO bookings (id, slot_id, phone_number, code, first_name, last_name) VALUES (DEFAULT, $1, $2, $3, $4, $5)"},
	KVP{K: "book/get", V: "SELECT * FROM bookings WHERE code=$1"},
}
