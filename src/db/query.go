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
	KVP{K: "book/remove", V: "DELETE FROM bookings WHERE code=$1"},
	KVP{K: "booking/get", V: "SELECT * FROM bookings"},
	KVP{K: "booking/getBySlotID", V: "SELECT * FROM bookings WHERE slot_id=$1"},
	KVP{K: "company/insert", V: "INSERT INTO company (id, name, adress, city, country, post_code, contact_firstname, contact_number, contact_lastname, verified, email, password, lon, lat, contact_email) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, DEFAULT, $9, $10, $11, $12, $13)"},
	KVP{K: "company/update/location", V: "UPDATE company SET name=$2, adress=$3, city=$4, country=$5, post_code=$6, lon=$7, lat=$8 WHERE id=$1 RETURNING *"},
	KVP{K: "company/update/contact", V: " UPDATE company SET contact_firstname=$2, contact_lastname=$3, contact_number=$4 WHERE id=$1 RETURNING *"},
	KVP{K: "company/update/password", V: "UPDATE company SET password=$2 WHERE id=$1 RETURNING *"},
	KVP{K: "company/login", V: "SELECT * FROM company WHERE email=$1"},
	KVP{K: "company/slot/get", V: "SELECT * FROM slots WHERE id=$1"},
	KVP{K: "company/distance", V: "SELECT * FROM company WHERE (lat >= $1 AND lat <= $2) AND (lon >= $3 AND lon <= $4) AND acos(sin($5) * sin(lat) + cos($5) * cos(lat) * cos(lon - $6)) <= $7"},
	KVP{K: "company/slot/getAll", V: "SELECT * FROM slots WHERE company_id=$1"},
	KVP{K: "company/slot/update", V: "UPDATE slots SET start_time=$2, end_time=$3, max=$4, day=$5 WHERE id=$1 RETURNING *"},
	KVP{K: "company/slot/add", V: "INSERT INTO slots (id, company_id, start_time, end_time, max, day) VALUES (DEFAULT, $1, $2, $3, $4, $5)"},
}
