package db

type KVP struct {
	K string
	V string
}

var queries = []KVP{
	KVP{K: "company/get", V: "SELECT * FROM company"},
	KVP{K: "company/getByID", V: "SELECT * FROM company WHERE id=$1"},
}
