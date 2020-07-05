package db

type KVP struct {
	K string
	V string
}

var queries = []KVP{
	KVP{K: "company/get", V: "SELECT * FROM company"},
	KVP{K: "company/getByID", V: "SELECT * FROM company WHERE id=$1"},
	KVP{K: "book/add", V: "INSERT INTO bookings (id, slot_id, phone_number, code, first_name, last_name, visitee, message, status) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8)"},
	KVP{K: "book/get", V: "SELECT * FROM bookings WHERE code=$1"},
	KVP{K: "book/remove", V: "DELETE FROM bookings WHERE code=$1"},
	KVP{K: "booking/getByPhone", V: "SELECT * FROM bookings WHERE phone_number=$1"},
	KVP{K: "booking/getByCompanyID", V: "SELECT b.id AS booking_id, * FROM bookings b JOIN slots s ON s.id=b.slot_id WHERE s.company_id=$1 AND (s.start_time >= $2 AND s.end_time <= $3) ORDER BY s.start_time"},
	KVP{K: "booking/getBySlotID", V: "SELECT * FROM bookings WHERE slot_id=$1"},
	KVP{K: "company/insert", V: "INSERT INTO company (id, name, adress, city, country, post_code, contact_firstname, contact_number, contact_lastname, verified, email, password, lon, lat, contact_email) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, DEFAULT, lower($9), $10, $11, $12, lower($13))"},
	KVP{K: "company/update/location", V: "UPDATE company SET name=$2, adress=$3, city=$4, country=$5, post_code=$6, lon=$7, lat=$8 WHERE id=$1 RETURNING *"},
	KVP{K: "company/update/contact", V: " UPDATE company SET contact_firstname=$2, contact_lastname=$3, contact_number=$4, contact_email=lower($5) WHERE id=$1 RETURNING *"},
	KVP{K: "company/update/password", V: "UPDATE company SET password=$2 WHERE id=$1 RETURNING *"},
	KVP{K: "company/login", V: "SELECT * FROM company WHERE lower(email)=lower($1)"},
	KVP{K: "company/slot/get", V: "SELECT * FROM slots WHERE id=$1"},
	KVP{K: "company/distance", V: "SELECT * FROM company WHERE (lat >= $1 AND lat <= $2) AND (lon >= $3 AND lon <= $4) AND acos(sin($5) * sin(lat) + cos($5) * cos(lat) * cos(lon - $6)) <= $7"},
	KVP{K: "company/slot/getAll", V: "SELECT * FROM slots WHERE company_id=$1 ORDER BY start_time"},
	KVP{K: "company/slot/update", V: "UPDATE slots SET start_time=$2, end_time=$3, max=$4 WHERE id=$1 RETURNING *"},
	KVP{K: "company/slot/add", V: "INSERT INTO slots (id, company_id, start_time, end_time, max) VALUES (DEFAULT, $1, $2, $3, $4)"},
	KVP{K: "company/slot/get/betweenTime", V: "select * from slots where start_time between $1 AND $2 AND company_id = $3 ORDER BY start_time"},
	KVP{K: "company/get/verified", V: "SELECT * FROM COMPANY WHERE verified=true"},
	KVP{K: "company/get/avgAvailability", V: `SELECT coalesce((sum(max) - sum(booked)) / sum(max)::float, 0)  as avg
										        FROM generate_series($2::timestamp, ($2::date + $3::integer)::timestamp, '1 day') t(day)
												LEFT JOIN slots s ON s.start_time::date=t.day::date AND s.company_id=$1
												GROUP BY t.day::date
												ORDER BY t.day::date`},
	KVP{K: "company/get/slotAvailability", V: `SELECT count(s.id) as avg
										        FROM generate_series($2::timestamp, ($2::date + $3::integer)::timestamp, '1 day') t(day)
												LEFT JOIN slots s ON s.start_time::date=t.day::date AND s.company_id=$1 AND s.booked < s.max
												GROUP BY t.day::date
												ORDER BY t.day::date`},
	KVP{K: "booking/update/status", V: "UPDATE bookings b SET status=$3 FROM slots s WHERE b.id=$2 AND s.company_id=$1 AND b.slot_id=s.id RETURNING b.*"},
	KVP{K: "company/search", V: `SELECT * FROM company c WHERE (($1::float4 IS NULL OR $2::float4 IS NULL) OR (lat >= $1 AND lat <= $2) AND (lon >= $3 AND lon <= $4) 
								AND acos(sin($5) * sin(lat) + cos($5) * cos(lat) * cos(lon - $6)) <= $7)
								AND ($8::varchar(50) IS NULL OR LOWER(c.name) LIKE LOWER(CONCAT('%', $8, '%')))`},
}
