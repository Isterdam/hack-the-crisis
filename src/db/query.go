package db

type KVP struct {
	K string
	V string
}

var queries = []KVP{
	KVP{K: "company/get", V: "SELECT * FROM company"},
	KVP{K: "company/getByID", V: "SELECT * FROM company WHERE id=$1"},
	KVP{K: "booking/get", V: "SELECT * FROM bookings"},
	KVP{K: "company/insert", V: "INSERT INTO company (id, name, adress, city, country, post_code, contact_firstname, contact_number, contact_lastname, verified, email, password) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, DEFAULT, $9, $10)"},
}
